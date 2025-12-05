//api入口
import { Router } from "express";
import testRouter from "./test.js"

const router = Router();

router.use("/test", testRouter)

export default router;