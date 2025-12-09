import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.get("/", async (req, res) => {
    try {
        let banner = await db.collection('banner').get()
        return res.json(ok(banner))
    } catch (e) { console.log(e) }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body;
        if (query._id) {
            let bannerObj = JSON.parse(JSON.stringify(query))
            delete bannerObj._id
            bannerObj.interval = Number(bannerObj.interval)
            await db.collection('banner').where({ _id: query._id }).update({
                ...bannerObj,
                updatedAt: new Date()
            })
            // console.log(res)
            return res.json(ok("更新成功"))
        } else {
            await db.collection("banner").add({
                detail: query.detail,
                interval: Number(query.interval),
                autoplay: query.autoplay,
                type: query.type,
                url: query.url,
                createdAt: new Date(),
                title: query.title
            })
            return res.json(ok("新增成功"))
        }
    } catch (e) { console.log(e) }
})

router.post("/del", async (req, res) => {
    try {
        const query = req.body
        if (query._id) {
            await db.collection('banner').where({ _id: query._id }).remove()
            return res.json(ok("删除成功"))
        }
    } catch (e) { console.log(e) }
})

export default router;