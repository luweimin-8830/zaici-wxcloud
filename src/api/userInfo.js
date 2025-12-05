//模版页面,新接口页面直接复制此页
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
                name: ''
            }
            const addRes = await db.collection('users').add(newUser)
            userObj = {...newUser,_id: addRes.id}
        }else{
            //老用户
            const existingUser = userQuery.data[0];
            const docId = existingUser._id;
            await db.collection('users').doc(docId).update({lastLogin:new Date()});
            userObj = { ...existingUser, lastLogin: new Date() };
        }

        res.json(ok(userObj))
    } catch (e) { console.error(e);next(e); }
})

export default router;