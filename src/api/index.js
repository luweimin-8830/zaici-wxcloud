//api入口
import { Router } from "express";
import { getDb } from "../../index";
import { getModels } from "../../index";
import test from "./test"

const router = Router();

router.use("/test",test)

export default router;