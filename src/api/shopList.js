import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.post("/", async (req, res, next) => {
    try {
        const query = req.body;
        let longitude = typeof query.longitude == "number" ? query.longitude : Number(query.longitude);
        let latitude = typeof query.latitude == "number" ? query.latitude : Number(query.latitude);
        let distance = typeof query.distance == "number" ? query.distance : Number(query.distance);
        const myLocation = new db.Geo.Point(longitude, latitude)
        if (distance >= 1000) {
            var SHOW_AS_KM = true
        } else {
            var SHOW_AS_KM = false
        }
        if (query.longitude) {
            let merchantList = await db.collection("shop_list")
                .aggregate()
                .geoNear({
                    distanceField: "distance", // 输出的每个记录中 distance 即是与给定点的距离
                    spherical: true,
                    near: myLocation,
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
                const shopIds = shopList.map(shop => shop._id)
                const onlineStats = await db.collection('online').aggregate()
                    .match({
                        shopId: _.in(shopIds),
                        status: '在线',
                        dueTime: _.gte(now)
                    })
                    .group({
                        _id: '$shopId',
                        total: { $sum: 1 }
                    })
                    .end();
                const countMap = {};
                if (onlineStats.data) {
                    onlineStats.data.forEach(item => {
                        countMap[item._id] = item.total;
                    });
                }
                // 4. 遍历店铺列表，将人数填进去 (内存匹配)
                shopList.forEach(shop => {
                    // 如果 map 里有值就取值，没有就是 0
                    shop.onlineCount = countMap[shop._id] || 0;
                });

                if (query.shopId != '') {
                    shopList = moveObjectToFirst(shopList, '_id', query.shopId);
                }
                merchantList.data = shopList
            }
            return res.json(ok(merchantList))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.error(e); next(e) }
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

router.post("/save", async (req, res) => {
    try {
        const query = req.body
        if (query.shopList) {
            query.shopList.createdAt = new Date()
            const shopList = await db.collection('shop_list').add(query.shopList);
            return res.json(ok(shopList))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/update", async (req, res) => {
    try {
        const query = req.body
        if (query.shopList._id) {
            let obj = query.shopList
            delete obj.createdAt
            delete obj.starttime
            delete obj.endtime
            let id = obj._id
            delete obj._id
            const shopList = await db.collection('shop_list').where({ _id: id }).update({
                ...obj,
                updatedAt: new Date()
            })
            return res.json(ok(shopList))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/del", async (req, res) => {
    try {
        const query = req.body
        if (query.id) {
            const shopList = await db.collection('shop_list').where({ _id: query.id }).remove();
            return res.json(ok("删除成功"))
        } else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;