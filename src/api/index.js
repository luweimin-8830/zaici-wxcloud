//api入口
import { Router } from "express";
import userInfoRouter from "./userInfo.js";
import shopListRouter from "./shopList.js";
import detailRecordRouter from "./detailRecord.js";
import onlineRouter from "./online.js";
import matchRouter from "./match.js";
import chatRouter from "./chat.js";

const router = Router();

router.use("/userInfo",userInfoRouter);
router.use("/shopList",shopListRouter);
router.use("/detailRecord",detailRecordRouter);
router.use("/online",onlineRouter);
router.use("/match",matchRouter);
router.use("/chat",chatRouter);

export default router;