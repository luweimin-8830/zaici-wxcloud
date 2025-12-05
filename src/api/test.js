//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";

const router = Router();

router.get("/", async (req, res) => {
    try {
        const data = "这是一个测试,hello."
        res.json(ok(data))
    } catch (e) { console.log(e) }
})

export default router;