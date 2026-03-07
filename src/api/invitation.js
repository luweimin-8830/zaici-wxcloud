import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getTcb } from "../util/tcb.js";
import moment from 'moment-timezone'
import axios from "axios";
const router = Router();
const db = getDb();
const tcb = getTcb();
const _ = db.command

router.post("/save", async (req, res) => {
    try {
        const { id, title, imageUrl, activity, status, posterUrl, requiredFields } = req.body;
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
            requiredFields: requiredFields || [],
            openId: OPENID,
            updatedAt: db.serverDate()
        };

        let invitationId = id;
        if (id) {
            // 更新逻辑
            await db.collection('invitation_demo').doc(id).update(data);
            return res.json(ok({ id, message: "邀请函更新成功" }));
        } else {
            // 新增逻辑
            data.qrcode = "";
            data.createdAt = db.serverDate();
            const result = await db.collection('invitation_demo').add(data);
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
                await db.collection('invitation_demo').doc(invitationId).update({
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
        const { id, title, imageUrl, activity, qrcode, status, posterUrl, requiredFields } = req.body;
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
        if (requiredFields !== undefined) updateData.requiredFields = requiredFields;

        await db.collection('invitation_demo').doc(id).update(updateData);

        res.json(ok({ message: "邀请函更新成功" }));
    } catch (e) {
        console.error("更新邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/join", async (req, res) => {
    try {
        const { invitationId, nickname, avatar, company, department } = req.body;
        const OPENID = req.headers["x-wx-openid"];

        if (!invitationId) {
            return res.json(fail(400, "缺少活动ID"));
        }

        // 1. 检查活动是否存在且在进行中
        const invRes = await db.collection('invitation_demo').doc(invitationId).get();
        const invitation = Array.isArray(invRes.data) ? invRes.data[0] : invRes.data;
        
        if (!invitation) {
            return res.json(fail(404, "活动不存在"));
        }
        if (invitation.status === "已停止") {
            return res.json(fail(400, "该活动报名已结束"));
        }

        // 2. 检查是否已经报名过
        const checkRes = await db.collection('inviter_demo').where({
            invitationId,
            openId: OPENID
        }).count();

        if (checkRes.total > 0) {
            return res.json(fail(400, "您已报名该活动，无需重复报名"));
        }

        // 3. 写入报名表
        const data = {
            invitationId,
            openId: OPENID,
            nickname: nickname || "匿名用户",
            avatar: avatar || "",
            company: company || "",
            department: department || "",
            createdAt: db.serverDate(),
            updatedAt: db.serverDate()
        };

        await db.collection('inviter_demo').add(data);

        // 4. 同步更新用户主表 (与 userInfo.js 字段对齐)
        const userUpdateData = {
            updatedAt: db.serverDate()
        };
        if (nickname) userUpdateData.name = nickname; // userInfo.js 中使用 name
        if (avatar) userUpdateData.avatar = avatar;
        if (company) userUpdateData.company = company;
        if (department) userUpdateData.department = department;

        await db.collection('users_demo').where({ openId: OPENID }).update(userUpdateData);

        // 5. 如果修改了昵称或头像，同步更新在线表 (对齐 userInfo.js 逻辑)
        if (nickname || avatar) {
            const onlineUpdate = {};
            if (nickname) onlineUpdate.name = nickname;
            if (avatar) onlineUpdate.avatar = avatar;
            await db.collection('online_demo').where({ openId: OPENID }).update(onlineUpdate);
        }

        res.json(ok({ message: "报名成功" }));
    } catch (e) {
        console.error("报名活动失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/joiners", async (req, res) => {
    try {
        const { invitationId } = req.body;
        if (!invitationId) {
            return res.json(fail(400, "缺少活动ID"));
        }

        const result = await db.collection('inviter_demo')
            .where({ invitationId })
            .orderBy('createdAt', 'desc')
            .limit(100)
            .get();

        res.json(ok(result.data));
    } catch (e) {
        console.error("获取报名列表失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/list", async (req, res) => {
    try {
        const { page = 1, limit = 10, keyword = "", activity, status } = req.body;
        const skip = (page - 1) * limit;
        
        let query = {};
        if (keyword) {
            query.title = db.RegExp({
                regexp: keyword,
                options: 'i',
            });
        }
        if (activity) {
            query.activity = activity;
        }
        if (status) {
            query.status = status;
        }

        const countRes = await db.collection('invitation_demo').where(query).count();
        const listRes = await db.collection('invitation_demo')
            .where(query)
            .orderBy('updatedAt', 'desc')
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

        const result = await db.collection('invitation_demo').doc(id).get();
        if (!result.data || (Array.isArray(result.data) && result.data.length === 0)) {
            return res.json(fail(404, "记录不存在"));
        }

        res.json(ok(Array.isArray(result.data) ? result.data[0] : result.data));
    } catch (e) {
        console.error("获取邀请函详情失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/checkJoin", async (req, res) => {
    try {
        const { invitationId } = req.body;
        const OPENID = req.headers["x-wx-openid"];

        if (!invitationId) {
            return res.json(fail(400, "缺少活动ID"));
        }

        const countRes = await db.collection('inviter_demo').where({
            invitationId,
            openId: OPENID
        }).count();

        res.json(ok({ isJoined: countRes.total > 0 }));
    } catch (e) {
        console.error("检查报名状态失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

router.post("/del", async (req, res) => {
    try {
        const { id } = req.body;
        if (!id) {
            return res.json(fail(400, "缺少记录ID"));
        }
        await db.collection('invitation_demo').doc(id).remove();
        res.json(ok({ message: "删除成功" }));
    } catch (e) {
        console.error("删除邀请函失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

// --- 抽奖相关接口 ---

/**
 * 保存/更新抽奖配置
 */
router.post("/lottery/save", async (req, res) => {
    try {
        const { id, invitationId, prizeName, winnerCount, winners, status } = req.body;
        const OPENID = req.headers["x-wx-openid"];

        if (!invitationId || !prizeName || !winnerCount) {
            return res.json(fail(400, "参数缺失: invitationId, prizeName, winnerCount 为必填项"));
        }

        const data = {
            invitationId,
            prizeName,
            winnerCount: parseInt(winnerCount),
            winners: winners || [], // 中奖人 openId 列表
            status: status || "进行中",
            creatorOpenId: OPENID,
            updatedAt: db.serverDate()
        };

        if (id) {
            await db.collection('lottery_demo').doc(id).update(data);
            return res.json(ok({ id, message: "抽奖更新成功" }));
        } else {
            data.createdAt = db.serverDate();
            const result = await db.collection('lottery_demo').add(data);
            return res.json(ok({ id: result.id, message: "抽奖创建成功" }));
        }
    } catch (e) {
        console.error("保存抽奖失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

/**
 * 获取活动下的抽奖列表
 */
router.post("/lottery/list", async (req, res) => {
    try {
        const { invitationId } = req.body;
        if (!invitationId) {
            return res.json(fail(400, "缺少活动ID"));
        }

        const result = await db.collection('lottery_demo')
            .where({ invitationId })
            .orderBy('createdAt', 'desc')
            .get();

        res.json(ok(result.data));
    } catch (e) {
        console.error("获取抽奖列表失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

/**
 * 获取抽奖详情
 */
router.post("/lottery/detail", async (req, res) => {
    try {
        const { id } = req.body;
        if (!id) {
            return res.json(fail(400, "缺少抽奖ID"));
        }

        const result = await db.collection('lottery_demo').doc(id).get();
        const data = Array.isArray(result.data) ? result.data[0] : result.data;
        
        if (!data) {
            return res.json(fail(404, "抽奖不存在"));
        }

        res.json(ok(data));
    } catch (e) {
        console.error("获取抽奖详情失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

/**
 * 删除抽奖
 */
router.post("/lottery/del", async (req, res) => {
    try {
        const { id } = req.body;
        if (!id) {
            return res.json(fail(400, "缺少抽奖ID"));
        }
        await db.collection('lottery_demo').doc(id).remove();
        res.json(ok({ message: "删除成功" }));
    } catch (e) {
        console.error("删除抽奖失败:", e);
        res.json(fail(500, "服务器内部错误"));
    }
});

export default router;
