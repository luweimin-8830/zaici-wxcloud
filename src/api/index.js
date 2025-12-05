//api入口
import { Router } from "express";
import testRouter from "./test.js"
import userInfoRouter from "./userInfo.js"

const router = Router();

router.use("/test", testRouter)
router.use("/userInfo",userInfoRouter)

export default router;