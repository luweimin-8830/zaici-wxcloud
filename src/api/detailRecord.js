import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.post("/get", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let detail_record = await db.collection("detail_record").where({ openId: query.openId, status: 1 }).get()
            let users = await db.collection('users').where({ openId: query.openId }).get()
            let info = {}
            if (users.data.length > 0 && detail_record.data.length > 0) {
                info = detail_record.data[0]
                info.avatar = users.data[0].avatar
                info.name = users.data[0].name
                return res.json(ok(info))
            } else {
                return res.json(ok("未填写今日的我"))
            }

        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let detail_record = await db.collection('detail_record').where({
                openId: query.openId,
                status: 1
            }).get()
            let detail = {
                openId: query.openId,
                plan: query.plan,
                mood: query.mood,
                style: query.style,
                description: query.description,
                seat: query.seat,
                image: query.image,
            }
            if (detail_record.data.length == 0) {
                let newDetail = {
                    ...detail,
                    status: 1,
                    createdAt: new Date()
                }
                await db.collection('detail_record').add(newDetail)
                return res.json(ok(newDetail))
            } else {
                await db.collection('detail_record').doc(detail_record.data[0]._id).update(detail)
                return res.json(ok(detail))
            }
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/saveSeat", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            await db.collection("detail_record").where({ openId: query.openId, status: 1 }).update({
                seat: query.seat
            })
            return res.json(ok("更新成功"))
        }else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;