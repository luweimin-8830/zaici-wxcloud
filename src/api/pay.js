//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/", async (req, res) => {
    try {
        const OPENID = req.headers['x-wx-openid'];
        let payment = {
            "body": "测试微信支付",
            "openid": OPENID,
            "out_trade_no": "1217752501201407033233368019",
            "spbill_create_ip": "127.0.0.1",
            "env_id": "prod-3g90nhycc15ce33f",
            "sub_mch_id": "1727232939",
            "total_fee": 1,
            "callback_type": 2,            
        }
        const wx_order_pay = await fetch("http://api.weixin.qq.com/_/pay/unifiedOrder", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payment)
        })
        console.log("wx payment result:",wx_order_pay)
        return res.json(ok(wx_order_pay))
    } catch (e) { console.log(e) }
})

export default router;