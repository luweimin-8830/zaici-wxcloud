import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels, getTcb } from "../util/tcb.js";
import moment from 'moment-timezone'
import axios from "axios";
const router = Router();
const db = getDb();
const tcb = getTcb();
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
        const invitationId = result.id;

        // 生成小程序码
        try {
            const qrRes = await axios.post("http://api.weixin.qq.com/wxa/getwxacode", {
                path: `pages/activityEntry/activityEntry?id=${invitationId}`,
                width: 430
            }, {
                responseType: 'arraybuffer'
            });

            // 检查是否返回错误信息
            if (qrRes.headers['content-type'] && qrRes.headers['content-type'].includes('application/json')) {
                const errData = JSON.parse(Buffer.from(qrRes.data).toString());
                console.error("微信生成小程序码接口报错:", errData);
            } else {
                const buffer = qrRes.data;
                
                // 上传到云存储
                const uploadRes = await tcb.uploadFile({
                    cloudPath: `qrcode/invitation_${invitationId}.png`,
                    fileContent: Buffer.from(buffer)
                });

                // 更新记录
                await db.collection('invitation').doc(invitationId).update({
                    qrcode: uploadRes.fileID,
                    updatedAt: db.serverDate()
                });
                
                return res.json(ok({ id: invitationId, qrcode: uploadRes.fileID, message: "邀请函保存成功" }));
            }
        } catch (qrErr) {
            console.error("小程序码生成流程异常:", qrErr);
        }
        
        res.json(ok({ id: invitationId, qrcode: "", message: "邀请函保存成功，小程序码生成失败" }));
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
