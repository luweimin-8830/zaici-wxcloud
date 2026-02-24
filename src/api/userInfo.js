import { Router } from "express";
import axios from "axios";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.post("/", async (req, res, next) => {
    try {
        const OPENID = req.headers["x-wx-openid"];
        if (!OPENID) { return res.json(fail(401, "未获取到openId,请确认")) }
        const userQuery = await db.collection('users').where({ openId: OPENID }).get()
        const userName = "用户" + OPENID.slice(-4)
        let userObj = {}
        if (userQuery.data.length === 0) {
            //新用户
            const newUser = {
                openId: OPENID,
                superLike: 0,
                beLike: 0,
                beSuperLike: 0,
                image: [],
                created: db.serverDate(),
                lastLogin: db.serverDate(),
                name: userName,
                avatar: "https://cloud1-1gth9cum37c9015c-1380861431.tcloudbaseapp.com/logo.png?sign=44d3bfbeb3cb05b7ff6fce3460a8bcdf&t=1765181560"
            }
            const addRes = await db.collection('users').add(newUser)
            userObj = { ...newUser, _id: addRes.id }
        } else {
            //老用户
            const existingUser = userQuery.data[0];
            const docId = existingUser._id;
            await db.collection('users').doc(docId).update({ lastLogin: new Date() });
            userObj = { ...existingUser, lastLogin: new Date() };
        }

        res.json(ok(userObj))
    } catch (e) { console.error(e); next(e); }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            await db.collection('users').where({ openId: query.openId })
                .update({
                    [query.key]: query.key == 'birthday' ? query.data ? new Date(query.data) : null : query.data,
                    updatedAt: new Date(),
                })
            if (query.key == 'avatar' || query.key == 'name') {
                await db.collection('online').where({ openId: query.openId })
                    .update({
                        [query.key]: query.data
                    })
            }
            return res.json(ok("更新成功"))
        } else {
            res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//变更super like次数
router.post("/superLike", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let user = await db.collection('users').where({ openId: query.openId }).get()
            if (user.data[0].superLike == 0) {
                return res.json(ok("superLike次数已用完"))
            } else if (user.data[0].superLike > 0) {
                user.data[0].superLike = user.data[0].superLike - 1
                await db.collection('users').where({ openId: query.openId }).update({
                    superLike: user.data[0].superLike
                })
                return res.json(ok("super Like次数-1"))
            }
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//发放super like次数
router.post("/addSuperLike", async (req, res) => {
    try {
        const query = req.body
        if (query.userList.length > 0) {
            for (let i = 0; i < query.userList.length; i++) {
                let user = await db.collection('users').where({ openId: query.userList[i].openId }).get()
                let total = user.data[0].superLike + query.total
                await db.collection('users').where({ openId: query.userList[i].openId }).update({
                    superLike: total
                })
            }
            return res.json(ok("发放成功"))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/addAdmin", async (req, res) => {
    try {
        const OPENID = req.headers["x-wx-openid"];
        let isAdmin = await db.collection('admins').where({ openId: OPENID }).get()
        if (isAdmin.data.length > 0) {
            return res.json(ok("已是管理员"))
        } else {
            await db.collection("admins").add({ openId: OPENID })
            return res.json(ok("添加管理员成功"))
        }
    } catch (e) { console.log(e) }
})

router.post("/getOtherUserInfo", async (req, res) => {
    try {
        const query = req.body;
        const viewerOpenId = query.openId || req.headers["x-wx-openid"];
        if (!query.id) {
            return res.json(fail(401, "参数错误"))
        }
        // 1) 查 online 记录
        const onlineDoc = await db.collection('online').doc(query.id).get()
        if (!onlineDoc.data || onlineDoc.data.length === 0) {
            return res.json(fail(404, "未找到用户在线记录"))
        }
        let otherUser = onlineDoc.data[0]

        // 2) 获取个人资料与用户信息
        const detailRes = await db.collection('detail_record').where({ openId: otherUser.openId, status: 1 }).get()
        const userRes = await db.collection('users').where({ openId: otherUser.openId }).get()
        otherUser.detailRecord = detailRes.data[0] || {}
        otherUser.userInfo = userRes.data[0] || {}

        // 3) 补齐展示用图片
        if (!otherUser.detailRecord.image || otherUser.detailRecord.image.length === 0) {
            if (otherUser.userInfo.avatar) {
                otherUser.detailRecord.image = [{ url: otherUser.userInfo.avatar }]
            } else {
                otherUser.detailRecord.image = []
            }
        }

        // 4) 计算匹配状态（可选）
        otherUser.state = 0
        otherUser.likeType = 0
        otherUser.channel = ''
        if (viewerOpenId) {
            const matchList = await db.collection("match").where({ openId1: viewerOpenId, openId2: otherUser.openId }).get()
            const matchList1 = await db.collection("match").where({ openId1: otherUser.openId, openId2: viewerOpenId }).get()
            if (matchList.data.length > 0 || matchList1.data.length > 0) {
                if ((matchList.data.length > 0 && matchList.data[0].status == 2) || (matchList1.data.length > 0 && matchList1.data[0].status == 2)) {
                    otherUser.state = 1
                    if (matchList.data.length > 0) {
                        otherUser.likeType = matchList.data[0].likeType
                        otherUser.channel = matchList.data[0].channel
                    } else {
                        otherUser.likeType = 0
                        otherUser.channel = matchList1.data[0].channel
                    }
                } else if (matchList.data.length > 0 && matchList.data[0].status == 1) {
                    otherUser.state = 0
                    otherUser.likeType = matchList.data[0].likeType
                    otherUser.channel = matchList.data[0].channel
                } else if (matchList1.data.length > 0) {
                    otherUser.state = 0
                    otherUser.likeType = 0
                    otherUser.channel = matchList1.data[0].channel
                }
            }
        }

        return res.json(ok(otherUser))
    } catch (e) {
        console.log(e)
        return res.json(fail(500, "获取用户信息失败"))
    }
})

router.post("/getInfo", async (req, res) => {
    try {
        const query = req.body;
        if (query.openId) {
            let data = {}
            const userInfo = await db.collection('users').where({ openId: query.openId }).get();
            const detail = await db.collection('detail_record').where({ openId: query.openId }).get();
            data.userInfo = userInfo;
            data.detailRecord = detail;
            res.json(ok(data))
        } else {
            res.json(error(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;