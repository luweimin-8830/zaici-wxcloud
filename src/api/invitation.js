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
        const { id, title, imageUrl, activity, status, posterUrl } = req.body;
        const OPENID = req.headers["x-wx-openid"];
        
        if (!title || !imageUrl || !activity) {
            return res.json(fail(400, "参数缺失"));
        }

        const data = {
            title,
            imageUrl,
            activity,
            status: status || "进行中",
            posterUrl: posterUrl || "",
            openId: OPENID,
            updatedAt: db.serverDate()
        };

        let invitationId = id;
        if (id) {
            // 更新逻辑
            await db.collection('invitation').doc(id).update(data);
            return res.json(ok({ id, message: "邀请函更新成功" }));
        } else {
            // 新增逻辑
            data.qrcode = "";
            data.createdAt = db.serverDate();
            const result = await db.collection('invitation').add(data);
            invitationId = result.id;
        }

        // 生成小程序码
        try {
            const qrRes = await axios.post("http://api.weixin.qq.com/wxa/getwxacode", {
                path: `pages/activityEntry/activityEntry?id=${invitationId}`,
                width: 430,
                env_version:"trial"
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
        const { id, title, imageUrl, activity, qrcode, status, posterUrl } = req.body;
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
        if (status !== undefined) updateData.status = status;
        if (posterUrl !== undefined) updateData.posterUrl = posterUrl;

        await db.collection('invitation').doc(id).update(updateData);

        res.json(ok({ message: "邀请函更新成功" }));
    } catch (e) {
        console.error("更新邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/list", async (req, res) => {
    try {
        const { page = 1, limit = 10, keyword = "" } = req.body;
        const skip = (page - 1) * limit;
        
        let query = {};
        if (keyword) {
            query.title = db.RegExp({
                regexp: keyword,
                options: 'i',
            });
        }

        const countRes = await db.collection('invitation').where(query).count();
        const listRes = await db.collection('invitation')
            .where(query)
            .orderBy('createdAt', 'desc')
            .skip(skip)
            .limit(limit)
            .get();

        res.json(ok({
            list: listRes.data,
            total: countRes.total,
            page,
            limit
        }));
    } catch (e) {
        console.error("获取邀请函列表失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/detail", async (req, res) => {
    try {
        const { id } = req.body;
        if (!id) {
            return res.json(fail(400, "缺少记录ID"));
        }

        const result = await db.collection('invitation').doc(id).get();
        if (!result.data || (Array.isArray(result.data) && result.data.length === 0)) {
            return res.json(fail(404, "记录不存在"));
        }

        res.json(ok(Array.isArray(result.data) ? result.data[0] : result.data));
    } catch (e) {
        console.error("获取邀请函详情失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/del", async (req, res) => {
    try {
        const { id } = req.body;
        if (!id) {
            return res.json(fail(400, "缺少记录ID"));
        }
        await db.collection('invitation').doc(id).remove();
        res.json(ok({ message: "删除成功" }));
    } catch (e) {
        console.error("删除邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

export default router;
