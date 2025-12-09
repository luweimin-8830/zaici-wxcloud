//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb ,getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/", async (req, res) => {
    try {
        const data = "这是一个测试,hello."
        res.json(ok(data))
    } catch (e) { console.log(e) }
})

export default router;