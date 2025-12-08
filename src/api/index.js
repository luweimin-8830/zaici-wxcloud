//api入口
import { Router } from "express";
import userInfoRouter from "./userInfo.js";
import shopListRouter from "./shopList.js";
import detailRecordRouter from "./detailRecord.js"

const router = Router();

router.use("/userInfo",userInfoRouter);
router.use("/shopList",shopListRouter);
router.use("/detailRecord",detailRecordRouter)

export default router;