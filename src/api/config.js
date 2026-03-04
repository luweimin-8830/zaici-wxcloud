//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb ,getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/saveDistance", async (req, res) => {
    try {
        const query = req.body
        await db.collection("config").doc("7dc2756e69a7bee40021341f7494b5b1").update({
            distance: query.distance
        })
        res.json(ok("更新成功"))
    } catch (e) { console.log(e) }
})

router.get("/getDistance", async (req, res) => {
    try {
        const distance = await db.collection("config").doc("7dc2756e69a7bee40021341f7494b5b1").get()
        res.json(ok(distance.data[0]))
    } catch (e) { console.log(e) }
})

export default router;