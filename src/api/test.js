//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/insert", async (req, res) => {
    try {
        const { data } = req.body;
        const result = await db.collection('test_demo').add({
            ...data,
            createdAt: new Date(),
            updatedAt: new Date()
        });
        res.json(ok(result));
    } catch (e) { 
        console.error(e);
        res.json(fail(500, "插入失败"));
    }
});

router.post("/update", async (req, res) => {
    try {
        const { id, data } = req.body;
        await db.collection('test_demo').doc(id).update({
            ...data,
            updatedAt: new Date()
        });
        res.json(ok("更新成功"));
    } catch (e) {
        console.error(e);
        res.json(fail(500, "更新失败"));
    }
});

export default router;