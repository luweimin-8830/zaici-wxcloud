import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb } from "../util/tcb.js";

const router = Router();
const db = getDb();
const _ = db.command;

router.get("/detail", async (req, res) => {
    try {
        const id = req.query._id || req.query.id;
        console.log("获取门店详情 ID:", id);
        if (id) {
            let shopData = null;
            
            // 方式1：doc().get() -> 改为使用 where().get() 以避免 ID 不存在时抛出异常
            const docRes = await db.collection('shop_list_demo').where({ _id: id }).get();
            if (docRes.data && docRes.data.length > 0) {
                shopData = docRes.data[0];
            }

            if (shopData) {
                console.log("找到门店数据:", shopData.shopname);
                return res.json(ok(shopData));
            } else {
                console.log("未找到门店数据, ID:", id);
                return res.json(fail(404, "未找到该门店"));
            }
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { 
        console.error("获取详情错误:", e); 
        res.json(fail(500, "获取详情失败")) 
    }
})

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
            let merchantList = await db.collection("shop_list_demo")
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
                const onlineStats = await db.collection('online_demo').aggregate()
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
    return arr;
}

router.post("/save", async (req, res) => {
    try {
        const query = req.body
        if (query.shopList) {
            const data = { ...query.shopList };
            // 清理无用或空字段，防止插入数据库时出错
            delete data._id; 
            data.createdAt = new Date();
            data.updatedAt = new Date();
            
            const result = await db.collection('shop_list_demo').add(data);
            return res.json(ok(result))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { 
        console.error("保存门店失败:", e);
        return res.json(fail(500, "保存门店失败"))
    }
})

router.post("/update", async (req, res) => {
    try {
        const query = req.body
        if (query.shopList._id) {
            let obj = query.shopList
            // 基础信息字段
            delete obj.createdAt
            delete obj.starttime
            delete obj.endtime
            // 相册相关字段 (albumName, albumStatus) 会通过 obj 自动更新
            let id = obj._id
            delete obj._id
            const shopList = await db.collection('shop_list_demo').where({ _id: id }).update({
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
            const shopList = await db.collection('shop_list_demo').where({ _id: query.id }).remove();
            return res.json(ok("删除成功"))
        } else {
            return res.json(fail(401,"参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/admin", async (req, res) => {
    try {
        let { page = 1, limit = 10, keyword = '' } = req.body;
        page = Number(page);
        limit = Number(limit);
        const skip = (page - 1) * limit;
        let query = {};
        if (keyword) {
            query = {
                shopname: {
                    $regex: keyword, 
                    $options: 'i' // 忽略大小写
                }
            };
        }
        const [listResult, countResult] = await Promise.all([
            db.collection("shop_list_demo")
                .where(query)                   
                .orderBy('create_time', 'desc') 
                .skip(skip)                     
                .limit(limit)                   
                .get(),
            
            db.collection("shop_list_demo")
                .where(query)                   
                .count()
        ]);
        res.json(ok({list:listResult.data,
            total: countResult.total,
            page:page,
            limit:limit
        }))
    } catch (e) { console.log(e) }
})

router.get("/moments", async (req, res) => {
    try {
        const { shopId, page = 1, limit = 10 } = req.query;
        if (!shopId) return res.json(fail(401, "缺少门店ID"));

        const skip = (Number(page) - 1) * Number(limit);
        
        // 聚合查询，关联用户信息
        const result = await db.collection('album_demo')
            .aggregate()
            .match({ shopId, status: 1 })
            .sort({ createdAt: -1 })
            .skip(skip)
            .limit(Number(limit))
            .lookup({
                from: 'users_demo', // 根据用户指正，表名为 users_demo
                localField: 'openId',
                foreignField: 'openId',
                as: 'userInfo'
            })
            .end();

        return res.json(ok(result.data));
    } catch (e) {
        console.error("获取动态失败:", e);
        return res.json(fail(500, "获取失败"));
    }
});

router.post("/publishMoment", async (req, res) => {
    try {
        const { shopId, openId, title, imageList } = req.body;
        if (!shopId || !openId || (!title && (!imageList || imageList.length === 0))) {
            return res.json(fail(401, "参数错误，内容或图片不能为空"));
        }

        const data = {
            shopId,
            openId,
            title,
            imageList,
            status: 1, // 默认正常状态
            createdAt: new Date(),
            updatedAt: new Date()
        };

        const result = await db.collection('album_demo').add(data);
        return res.json(ok(result));
    } catch (e) {
        console.error("发布动态失败:", e);
        return res.json(fail(500, "发布失败"));
    }
});

export default router;