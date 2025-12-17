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
        const OPENID = req.body.openId
        if (OPENID) {
            await db.collection('information_monitor').where({ openId: OPENID }).remove()
            return res.json(ok("删除成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//文本检测
router.post("/check", async (req, res) => {
    try {
        const query = req.body
        console.log(query)
        const wx_msg_check = await fetch("http://api.weixin.qq.com/wxa/msg_sec_check", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "content": query.content,
                "version": 2,
                "scene": query.scene,
                "openid": query.openId, 
            })
        })
        const data = wx_msg_check.json();
        if (data.errcode && data.errcode !== 0) {
            return res.json(fail(500,data));
        }
        return res.json(ok(data))
    } catch (e) { console.log(e) }
})


export default router;