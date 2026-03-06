import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.get("/", async (req, res) => {
    try {
        let banner = await db.collection('banner_demo').get()
        return res.json(ok(banner.data))
    } catch (e) { console.log(e) }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body;
        const bannerDate = query.banner;
        if (query.banner._id) {
            let updateData = {...bannerDate};
            delete updateData._id
            if (updateData.interval) {
                updateData.interval = Number(updateData.interval);
            }
            await db.collection('banner_demo').doc( bannerDate._id ).update({
                ...updateData,
                updatedAt: new Date()
            })
            return res.json(ok("更新成功"))
        } else {
            await db.collection("banner_demo").add({
                ...query.banner,
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