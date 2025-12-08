//模版页面,新接口页面直接复制此页
import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/", async (req, res, next) => {
    try {
        let longitude = typeof req.longitude == "number" ? req.longitude : Number(req.longitude);
        let latitude = typeof req.latitude == "number" ? req.latitude : Number(req.latitude);
        let distance = typeof req.distance == "number" ? req.distance : Number(req.distance);
        if (distance >= 1000) {
            var SHOW_AS_KM = true
        } else {
            var SHOW_AS_KM = false
        }
        if (req.body.longitude) {
            let merchantList = await db.collection("shop_list")
                .aggregate()
                .geoNear({
                    distanceField: "distance", // 输出的每个记录中 distance 即是与给定点的距离
                    spherical: true,
                    near: new db.Geo.Point(longitude, latitude),
                    minDistance: 0,
                    maxDistance: distance,
                    distanceMultiplier: SHOW_AS_KM ? 0.001 : 1,
                    key: "location", // 若只有 location 一个地理位置索引的字段，则不需填
                    includeLocs: "location", // 若只有 location 一个是地理位置，则不需填
                })
                .end()
            if (merchantList.data.length > 0) {
                let shopList = merchantList.data
                const now = new Date().getTime()
                for (let i = 0; i < shopList.length; i++) {
                    let count = await db.collection('online').where({
                        shopId: shopList[i]._id,
                        status: '在线',
                        dueTime: _.gte(now)
                    }).count()
                    shopList[i].onlineCount = count.total
                }
                if (req.shopId != '') {
                    shopList = moveObjectToFirst(shopList, '_id', req.shopId);
                }
                merchantList.data = shopList
            }
            return merchantList
        } else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.error(e);next(e) }
})

function moveObjectToFirst(arr, key, value) {
    // 查找匹配对象的索引
    const index = arr.findIndex(item => item[key] === value);

    // 如果找到匹配的对象
    if (index !== -1) {
        // 从数组中移除该对象
        const [item] = arr.splice(index, 1);
        // 将对象插入到数组开头
        arr.unshift(item);
    }
    console.log(arr)
    return arr;
}

export default router;