import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";
import moment from 'moment-timezone'
const router = Router();
const db = getDb();
const _ = db.command

router.post("/save", async (req, res) => {
    try {
        const { title, imageUrl, activity } = req.body;
        const OPENID = req.headers["x-wx-openid"];
        
        if (!title || !imageUrl || !activity) {
            return res.json(fail(400, "参数缺失"));
        }

        const data = {
            title,
            imageUrl,
            activity,
            qrcode: "", // 留空
            openId: OPENID,
            createdAt: db.serverDate(),
            updatedAt: db.serverDate()
        };

        const result = await db.collection('invitation').add(data);
        
        res.json(ok({ id: result.id, message: "邀请函保存成功" }));
    } catch (e) {
        console.error("保存邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/update", async (req, res) => {
    try {
        const { id, title, imageUrl, activity, qrcode } = req.body;
        const OPENID = req.headers["x-wx-openid"];

        if (!id) {
            return res.json(fail(400, "缺少记录ID"));
        }

        const updateData = {
            updatedAt: db.serverDate()
        };

        if (title !== undefined) updateData.title = title;
        if (imageUrl !== undefined) updateData.imageUrl = imageUrl;
        if (activity !== undefined) updateData.activity = activity;
        if (qrcode !== undefined) updateData.qrcode = qrcode;

        await db.collection('invitation').doc(id).update(updateData);

        res.json(ok({ message: "邀请函更新成功" }));
    } catch (e) {
        console.error("更新邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

export default router;
