import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.get("/", async (req, res) => {
    try {
        let banner = await db.collection('banner_demo').get()
        return res.json(ok(banner.data))
    } catch (e) { console.log(e) }
})

router.get("/detail", async (req, res) => {
    try {
        const id = req.query.id;
        if (!id) return res.json(fail(401, "缺少参数ID"));
        const banner = await db.collection('banner_demo').where({ _id: id }).get();
        if (banner.data && banner.data.length > 0) {
            return res.json(ok(banner.data[0]));
        } else {
            return res.json(fail(404, "未找到该广告详情"));
        }
    } catch (e) {
        console.log(e);
        return res.json(fail(500, "服务器错误"));
    }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body;
        const bannerData = query.banner;
        
        // 确保 type 和 url 存入 detail 字段，采用增量赋值方式防止覆盖原有其他自定义字段
        if (!bannerData.detail) bannerData.detail = {};
        if (!bannerData.detail.type) bannerData.detail.type = bannerData.type || 'image';
        if (!bannerData.detail.url) bannerData.detail.url = bannerData.url || '';

        if (bannerData._id) {
            let updateData = {...bannerData};
            delete updateData._id
            if (updateData.interval) {
                updateData.interval = Number(updateData.interval);
            }
            await db.collection('banner_demo').doc( bannerData._id ).update({
                ...updateData,
                updatedAt: new Date()
            })
            return res.json(ok("更新成功"))
        } else {
            await db.collection("banner_demo").add({
                ...bannerData,
                createdAt: new Date(),
            })
            return res.json(ok("新增成功"))
        }
    } catch (e) { 
        console.log(e);
        return res.json(fail("服务器错误"));
    }
})

router.post("/del", async (req, res) => {
    try {
        const query = req.body
        if (query._id) {
            await db.collection('banner_demo').where({ _id: query._id }).remove()
            return res.json(ok("删除成功"))
        }
    } catch (e) { console.log(e) }
})

export default router;