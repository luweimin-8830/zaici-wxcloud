import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/get", async (req, res) => {
    try {
        const { channel, length, cont } = req.body;
        if (!channel) return res.json(fail(400, "缺少频道ID"));

        const list = await db.collection('new_chat_history_demo')
            .where({ channelId: channel })
            .orderBy('timestamp', 'desc')
            .skip(length || 0)
            .limit(cont || 20)
            .get();
            
        return res.json(ok(list));
    } catch (e) {
        console.error("获取聊天记录失败:", e);
        res.json(fail(500, "服务器错误"));
    }
});

// 保存消息到数据库
router.post("/send", async (req, res) => {
    try {
        const { channelId, senderOpenID, receiverOpenID, messageContent } = req.body;
        if (!channelId || !senderOpenID || !receiverOpenID) {
            return res.json(fail(400, "参数错误"));
        }

        const data = {
            channelId,
            senderOpenID,
            receiverOpenID,
            messageContent,
            timestamp: messageContent.id || Date.now(),
            createdAt: new Date()
        };

        const result = await db.collection('new_chat_history_demo').add(data);
        return res.json(ok(result));
    } catch (e) {
        console.error("保存消息失败:", e);
        res.json(fail(500, "服务器错误"));
    }
});

router.post("/update", async (req, res) => {
    try {
        const query = req.body
        if (query.channelId) {
            const chat = await db.collection('new_chat_history_demo')
                .where(_.and([
                    { channelId: query.channelId },
                    { receiverOpenID: query.openId }
                ]))
                .update({ 'messageContent.state': query.state })
            return res.json(ok("更新成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/updateState", async (req, res) => {
    try {
        const query = req.body
        if (query.channelId) {
            const chat = await db.collection('new_chat_history_demo')
                .where(_.and([
                    { channelId: query.channelId },
                    { "messageContent.id": query.id }
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
            await db.collection('block_list_demo').add({ openId: query.openId, blockId: query.blockId, createdAt: new Date() })
            return res.json(ok("屏蔽成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//删除提示
router.post("/delInfo", async (req, res) => {
    try {
        const OPENID = req.headers["x-wx-openid"];
        if (OPENID) {
            await db.collection('information_monitor_demo').where({ openId: OPENID }).remove()
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
        const data = await wx_msg_check.json();
        console.log(`检测违规内容 ${data}`)
        if (data.errcode && data.errcode !== 0) {
            return res.json(fail(500,data));
        }
        return res.json(ok(data))
    } catch (e) { console.log(e) }
})


export default router;