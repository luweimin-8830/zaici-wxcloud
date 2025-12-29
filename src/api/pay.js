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
        let orderNo = Date.now() + Math.random().toString().substr(2, 5)
        let payment = {
            "openid": OPENID,
            "body": "测试微信支付",
            "trade_type": "JSAPI",
            "out_trade_no": orderNo,
            "spbill_create_ip": "127.0.0.1",
            "env_id": "prod-3g90nhycc15ce33f",
            "sub_mch_id": "1727232939",
            "total_fee": 1,
            "callback_type": 2,
            "container": {
                "service": "express-gqsx",
                "path": "/payCallback"
            }
        }
        const wx_order_pay = await fetch("http://api.weixin.qq.com/_/pay/unifiedOrder", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payment)
        })
        console.log("wx payment result:", wx_order_pay)
        return res.json(ok(wx_order_pay))
    } catch (e) { console.log(e) }
})

export default router;