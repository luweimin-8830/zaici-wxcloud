import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";
import moment from 'moment-timezone'
const router = Router();
const db = getDb();
const _ = db.command

router.post("/status", async (req, res) => {
    try {
        const query = req.body;

        if (query.openId) {
            let user = await db.collection('online').where({
                openId: query.openId,
                dueTime: _.gte(Date.now())
            }).get()
            if (user.data.length > 0) {
                let position = user.data[0]
                return res.json(ok({ position: position, status: "已出门" }))
            } else {
                return res.json(ok({ status: "未出门" }))
            }
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/near", async (req, res) => {
    try {
        const query = req.body;
        if (!query.shopId || query.shopId == '') {
            let userList = []
            let longitude = typeof query.longitude == "number" ? query.longitude : Number(query.longitude);
            let latitude = typeof query.latitude == "number" ? query.latitude : Number(query.latitude);
            const now = new Date().getTime()
            let list = await db.collection('online').aggregate()
                .geoNear({
                    distanceField: "distance", // 输出的每个记录中 distance 即是与给定点的距离
                    spherical: true,
                    near: new db.Geo.Point(longitude, latitude),
                    minDistance: 0,
                    maxDistance: 10000,
                    distanceMultiplier: 1,
                    key: "location", // 若只有 location 一个地理位置索引的字段，则不需填
                    includeLocs: "location", // 若只有 location 一个是地理位置，则不需填
                }).match({
                    dueTime: _.gte(now),
                    openId: _.neq(query.openId),
                    flag: 2
                })
                .end()
            let blockList = await db.collection('block_list').where({ openId: query.openId }).get()
            for (let i = 0; i < list.data.length; i++) {
                let matchList = await db.collection("match").where({ openId1: query.openId, openId2: list.data[i].openId }).get()
                let matchList1 = await db.collection("match").where({ openId1: list.data[i].openId, openId2: query.openId }).get()
                if (matchList.data.length > 0 || matchList1.data.length > 0) {
                    if ((matchList.data.length > 0 && matchList.data[0].status == 2) || (matchList1.data.length > 0 && matchList1.data[0].status == 2)) {
                        list.data[i].state = 1
                        if (matchList.data.length > 0) {
                            list.data[i].likeType = matchList.data[0].likeType
                            list.data[i].channel = matchList.data[0].channel
                        } else {
                            list.data[i].likeType = 0
                            list.data[i].channel = matchList1.data[0].channel
                        }

                    } else if (matchList.data.length > 0 && matchList.data[0].status == 1) {
                        list.data[i].state = 0
                        list.data[i].likeType = matchList.data[0].likeType
                        list.data[i].channel = matchList.data[0].channel
                    } else if (matchList1.data.length > 0) {
                        list.data[i].state = 0
                        list.data[i].likeType = 0
                        list.data[i].channel = matchList1.data[0].channel
                    } else {
                        list.data[i].state = 0
                        list.data[i].likeType = 0
                        list.data[i].channel = ''
                    }
                    if (list.data[i].likeType == 2 && matchList.data.length > 0) {
                        list.data[i].score = 100
                    }
                } else {
                    list.data[i].state = 0
                    list.data[i].likeType = 0
                    list.data[i].channel = ''
                }
                let detail = await db.collection('detail_record').where({ openId: list.data[i].openId, status: 1 }).get()
                let userInfo = await db.collection('users').where({ openid: list.data[i].openId }).get()
                if (userInfo.data.length > 0 && userInfo.data[0].image.length > 0) {
                    list.data[i].score = list.data[i].score + 35
                }
                if (detail.data.length > 0) {
                    let mergedObj = { ...list.data[i], ...detail.data[0] };
                    if (detail.data[0].image.length > 0) {
                        list.data[i].score = list.data[i].score + 35
                    }
                    if (detail.data[0].image.length == 0) {

                        if (userInfo.data.length > 0 && userInfo.data[0].avatar) {
                            mergedObj.image = [{ url: userInfo.data[0].avatar }]
                        } else {
                            mergedObj.image = []
                        }
                    }
                    list.data[i] = mergedObj
                }
                if (blockList.data.length > 0) {
                    let flag = true
                    blockList.data.forEach((e) => {
                        if (e.blockId == list.data[i].openId) {
                            flag = false
                        }
                    })
                    if (flag) {
                        userList.push(list.data[i])
                    }
                } else {
                    userList.push(list.data[i])
                }
            }
            list.data = sort(userList)
            return res.json(ok(list))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/shop", async (req, res) => {
    try {
        const query = req.body;
        if (query.shopId) {
            const now = new Date().getTime()
            let userList = []
            let list = await db.collection('online').where({
                shopId: query.shopId,
                status: '在线',
                dueTime: _.gte(now),
                openId: _.neq(query.openId),
                flag: 1
            }).get()
            let blockList = await db.collection('block_list').where({ openId: query.openId }).get()
            for (let i = 0; i < list.data.length; i++) {
                list.data[i].score = 0
                let matchList = await db.collection("match").where({ openId1: query.openId, openId2: list.data[i].openId }).get()
                let matchList1 = await db.collection("match").where({ openId1: list.data[i].openId, openId2: query.openId }).get()
                if (matchList.data.length > 0 || matchList1.data.length > 0) {
                    if ((matchList.data.length > 0 && matchList.data[0].status == 2) || (matchList1.data.length > 0 && matchList1.data[0].status == 2)) {
                        list.data[i].state = 1

                        if (matchList.data.length > 0) {
                            list.data[i].likeType = matchList.data[0].likeType
                            list.data[i].channel = matchList.data[0].channel
                        } else {
                            list.data[i].likeType = 0
                            list.data[i].channel = matchList1.data[0].channel
                        }

                    } else if (matchList.data.length > 0 && matchList.data[0].status == 1) {
                        list.data[i].state = 0
                        list.data[i].likeType = matchList.data[0].likeType
                        list.data[i].channel = matchList.data[0].channel
                    } else if (matchList1.data.length > 0) {
                        list.data[i].state = 0
                        list.data[i].likeType = 0
                        list.data[i].channel = matchList1.data[0].channel
                    } else {
                        list.data[i].state = 0
                        list.data[i].likeType = 0
                        list.data[i].channel = ''
                    }
                    if (list.data[i].likeType == 2 && matchList.data.length > 0) {
                        list.data[i].score = 100
                    }
                } else {
                    list.data[i].state = 0
                    list.data[i].likeType = 0
                    list.data[i].channel = ''
                }

                let detail = await db.collection('detail_record').where({ openId: list.data[i].openId, status: 1 }).get()
                let userInfo = await db.collection('users').where({ openId: list.data[i].openId }).get()
                list.data[i].detailRecord = detail.data[0] || {}
                list.data[i].userInfo = userInfo.data[0] || {}
                if (userInfo.data.length > 0 && userInfo.data[0].image.length > 0) {
                    list.data[i].score = list.data[i].score + 35
                }
                if (detail.data.length > 0) {
                    let mergedObj = { ...list.data[i], ...detail.data[0] };
                    if (detail.data[0].image.length > 0) {
                        list.data[i].score = list.data[i].score + 35
                    }
                    if (detail.data[0].image.length == 0) {
                        if (userInfo.data.length > 0 && userInfo.data[0].avatar) {
                            mergedObj.image = [{ url: userInfo.data[0].avatar }]
                        } else {
                            mergedObj.image = []
                        }
                    }
                    list.data[i] = mergedObj
                }

                if (blockList.data.length > 0) {
                    let flag = true
                    blockList.data.forEach((e) => {
                        if (e.blockId == list.data[i].openId) {
                            flag = false
                        }
                    })
                    if (flag) {
                        userList.push(list.data[i])
                    }
                } else {
                    userList.push(list.data[i])
                }
            }

            list.data = sort(userList)
            // list.data = userList
            return res.json(ok(list))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

router.post("/save", async (req, res) => {
    try {
        const query = req.body
        let user = await db.collection('users').where({
            openId: query.openId,
        }).get()
        let name = ''
        let avatar = ''
        if (user.data.length > 0) {
            avatar = user.data[0].avatar
            name = user.data[0].name
        }
        let now = new Date()
        now.setDate(now.getDate() + 1); // 增加 1 天（UTC）
        now.setHours(5, 0, 0, 0); // 05:00:00 UTC// 强制设置为本地时间 05:00:00.000
        const due = new Date(now).getTime()
        let location = {}
        if (!query.location) {
            location = {
                coordinates: [0, 0],
                type: 'Point'
            }
        } else {
            location = query.location
        }
        // 把之前的在线记录设为结束
        await db.collection('online').where({
            openId: query.openId,
            dueTime: _.gt(Date.now())
        }).update({
            dueTime: Date.now()
        })

        // 新增一条进店记录
        await db.collection('online').add({
            openId: query.openId,
            shopId: query.shopId,
            shopName: query.shopName,
            name: name,
            location: location,
            avatar: avatar,
            status: '在线',
            createdAt: new Date(),
            updatedAt: new Date(),
            dueTime: due,
            flag: query.flag
        })
        let online = await db.collection('online').where({
            openId: query.openId,
            dueTime: db.command.gte(Date.now())
        }).get()

        return res.json(ok(online))
    } catch (e) { console.log(e) }
})

router.post("/update", async (req, res) => {
    try {
        const query = req.body
        if (query.id) {
            const now = new Date()
            let online = await db.collection('online').doc(query.id).update({
                dueTime: now
            })
            const data = { code: 0, online: online, message: "更新成功" }
            return res.json(ok(data))
        }else{
            return res.json(fail(401,"参数错误"))
        }

    } catch (e) { console.log(e) }
})

router.post("/history", async (req, res) => {
    try {
        const { shopId, openId } = req.body;
        if (!shopId) return res.json(fail(401, "参数错误: 缺失 shopId"));

        const now = new Date().getTime();
        const _ = db.command;

        // 1. 获取该门店已离开的用户记录，排除当前用户自己
        // 按 dueTime 倒序排序，即最晚离开的人排在最前面
        let listResult = await db.collection('online')
            .where({
                shopId: shopId,
                dueTime: _.lte(now),
                openId: _.neq(openId)
            })
            .orderBy('dueTime', 'desc')
            .limit(100)
            .get();

        let rawList = listResult.data;
        if (rawList.length === 0) return res.json(ok({ data: [] }));

        // 2. 按 openId 去重
        // 历史记录中，同一个用户可能多次出入同一门店，这里只保留其最近的一次足迹
        const uniqueRawList = [];
        const seenOpenIds = new Set();
        for (const item of rawList) {
            if (!seenOpenIds.has(item.openId)) {
                uniqueRawList.push(item);
                seenOpenIds.add(item.openId);
            }
        }

        // 3. 并行获取关联信息（性能提升：使用 Promise.all 替代循环 await）
        // 仅处理前 20 条，保证响应速度
        const processedList = await Promise.all(uniqueRawList.slice(0, 20).map(async (item) => {
            const targetOpenId = item.openId;

            // 并行查询：用户信息、今日的我、匹配/喜欢状态
            const [userRes, detailRes, matchRes] = await Promise.all([
                db.collection('users').where({ openId: targetOpenId }).get(),
                db.collection('detail_record').where({ openId: targetOpenId, status: 1 }).get(),
                db.collection('match').where(
                    _.or([
                        { openId1: openId, openId2: targetOpenId },
                        { openId1: targetOpenId, openId2: openId }
                    ])
                ).get()
            ]);

            const userInfo = userRes.data[0] || {};
            const detail = detailRes.data[0] || {};
            const matchInfo = matchRes.data[0] || null;

            // 处理匹配状态 (state) 和 喜欢类型 (likeType)
            let likeType = 0;
            let state = 0;
            if (matchInfo) {
                state = matchInfo.status === 2 ? 1 : 0;
                if (matchInfo.openId1 === openId) {
                    likeType = matchInfo.likeType || 0;
                }
            }

            // 4. 返回符合 nearList.vue 模板结构的完整对象
            return {
                ...item,
                name: userInfo.name || '在此用户',
                avatar: userInfo.avatar || '',
                mood: detail.mood || '今日心情',
                plan: detail.plan || '今日安排',
                likeType: likeType,
                state: state,
                detailRecord: {
                    ...detail,
                    image: (detail.image && detail.image.length > 0)
                        ? detail.image
                        : (userInfo.avatar ? [{ url: userInfo.avatar }] : [])
                }
            };
        }));

        return res.json(ok({ data: processedList }));

    } catch (e) {
        console.error("History API Error:", e);
        return res.json(fail(500, "服务器内部错误"));
    }
})

export default router;

export function sort(data) {
    if (data.length > 0) {
        data.forEach(item => {
            if (item.updatedAt) {
                const date1 = new Date(item.updatedAt);
                const date2 = new Date();
                // 验证日期是否有效
                if (isNaN(date1.getTime()) || isNaN(date2.getTime())) {
                    throw new Error('无效的日期格式');
                }
                // 计算时间差（毫秒）
                const timeDiff = date2.getTime() - date1.getTime();
                // 转换为分钟
                const minutesDiff = timeDiff / (1000 * 60);
                if (minutesDiff <= 10) {
                    item.score = item.score + 30
                } else if (minutesDiff > 10 && minutesDiff <= 30) {
                    item.score = item.score + 25
                } else if (minutesDiff > 30 && minutesDiff <= 60) {
                    item.score = item.score + 20
                } else if (minutesDiff > 60 && minutesDiff <= 120) {
                    item.score = item.score + 10
                }
            }
        })
        let key = 'score'
        return [...data].sort((a, b) => b.score - a.score)
    } else {
        return data
    }
}