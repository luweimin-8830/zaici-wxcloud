import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/get", async (req, res) => {
    try {
        query = req.body
        if (query.channel) {
            let list = await db.collection('new_chat_history').orderBy('timestamp', 'desc').where({
                channelId: query.channel
            }).skip(query.length).limit(query.cont).get()
            return res.json(ok(list))
        } else {
            res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/update", async (req, res) => {
    try {
        const query = req.body
        if (query.channelId) {
            const chat = await db.collection('new_chat_history')
                .where(_.and([
                    { channelId: query.channelId },
                    { 'messageContent.id': query._id }
                ]))
                .update({ 'messageContent.state': query.state })
            return res.json(ok("更新成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//屏蔽用户
router.post("/saveBlock", async (req, res) => {
    try {
        const query = req.body
        if (query.openId && query.blockId) {
            await db.collection('block_list').add({ openId: query.openId, blockId: query.blockId, createdAt: new Date() })
            return res.json(ok("屏蔽成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//删除提示
router.post("/delInfo", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let res = await db.collection('information_monitor').where({ openId: query.openId }).remove()
            return res.json(ok("删除成功"))
        } else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;