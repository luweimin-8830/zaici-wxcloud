import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();

router.post("/get", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let res = await db.collection("detail_record").where({ openId: query.openId, status: 1 }).get()
            let res1 = await db.collection('users').where({ openid: query.openId }).get()
            let info = {}
            if (res1.data.length > 0 && res.data.length > 0) {
                info = res.data[0]
                info.avatar = res1.data[0].avatar
                info.name = res1.data[0].name
                return res.json(ok(info))
            } else {
                return res.json(ok("未填写今日的我"))
            }

        } else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;