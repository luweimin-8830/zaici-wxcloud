import { Router } from "express";
import { ok, fail } from "../response.js";
import { getDb, getModels } from "../util/tcb.js";
import { sort } from "./online.js";

const router = Router();
const db = getDb();
const _ = db.command;

//获取匹配
router.post("/get", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            let user = await db.collection('match').where(_.and([
                // AND 的第一个条件：status为1
                { status: query.status },
                // AND 的第二个条件：字段1或字段2包含'要查询的值'
                _.or([
                    { openId1: query.openId },
                    { openId2: query.openId }
                ])
            ])).get()
            console.log(user)
            let list = []
            let avatar = ''
            let likeCount = await db.collection('match').where({ status: 1, openId2: query.openId }).get()
            if (likeCount.data.length > 0) {
                let likep = likeCount.data[0]
                if (likep.openId1 == query.openId) {
                    let r = await db.collection('users').where({ openId: likep.openId2 }).get()
                    if (r.data.length > 0) {
                        avatar = r.data[0].avatar
                    }
                } else if (likep.openId2 == query.openId) {
                    let r = await db.collection('users').where({ openId: likep.openId1 }).get()
                    if (r.data.length > 0) {
                        avatar = r.data[0].avatar
                    }
                }
            }
            let likeObj = {
                count: likeCount.data.length,
                avatar: likeCount.data.length == 0 ? 'https://cloudbase-3gn2elwa3387b385-1364843451.tcloudbaseapp.com/user-search-fill.png?sign=fd684b2ec65ae33bd702f70edc25997f&t=1757054747' : avatar,
                name: '喜欢我的人',
                content: 1,
                contentTime: new Date().getTime(),
            }
            list.push(likeObj)
            if (user.data.length > 0) {
                list = list.concat(user.data)
                for (let i = 1; i < list.length; i++) {
                    let id = ""
                    if (list[i].openId1 != query.openId) {
                        id = list[i].openId1
                    } else if (list[i].openId2 != query.openId) {
                        id = list[i].openId2
                    } else {
                        id = query.openId
                    }
                    let person = await db.collection('users').where({
                        openId: id
                    }).get()
                    let detail = await db.collection('detail_record').where({ openId: id, status: 1 }).get()
                    list[i].userInfo = person.data[0] || {}
                    list[i].detailRecord = detail.data[0] || {}
                    
                    if (person.data.length > 0) {
                        list[i].avatar = person.data[0].avatar
                        list[i].name = person.data[0].name == null ? "" : person.data[0].name
                    }
                    let chatContent = await db.collection('new_chat_history').where({ channelId: list[i].channel }).orderBy('timestamp', 'desc').limit(5).get();
                    if (chatContent.length != 0) {
                        list[i].content = chatContent.data
                        list[i].contentTime = chatContent.data[0] == null ? '' : chatContent.data[0].messageContent.id
                    }
                    let chatList = await db.collection('new_chat_history').where({ channelId: list[i].channel, senderOpenID: id }).get()
                    let wRead = 0
                    if (chatList.data.length > 0) {
                        chatList.data.forEach(item => {
                            if (item.messageContent.state == 0) {
                                wRead++
                            }
                        })
                    }
                    list[i].wRead = wRead
                }
                return res.json(ok(list))
            } else {
                return res.json(ok(list))
            }
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

//增加/变更匹配状态
router.post("/add", async (req, res) => {
    try {
        const query = req.body
        let openId1 = query.openId1
        let openId2 = query.openId2
        let operation = query.operation
        let channel = query.channel
        let matchList = await db.collection("match").where({ openId1: openId1, openId2: openId2 }).get()
        let matchList1 = await db.collection("match").where({ openId1: openId2, openId2: openId1 }).get()
        let status = ''
        if (matchList.data.length > 0) {
            if (operation == 1) {
                await db.collection('match').where({ channel: query.channel })
                    .update({
                        status: 1,
                        likeType: query.likeType,
                        updatedAt: new Date(),
                    })
                status = 1
                if (query.likeType == 2) {
                    let user = await db.collection("users").where({ openId: query.openId2 }).get()
                    changeData(query, user)
                }
            } else if (operation == 0) {
                await db.collection('match').doc(matchList.data[0]._id)
                    .update({
                        status: 3,
                        updatedAt: new Date(),
                    })
                status = 3
            }
        } else if (matchList1.data.length > 0) {
            if (operation == 1) {
                await db.collection('match').doc(matchList1.data[0]._id)
                    .update({
                        status: 2,
                        updatedAt: new Date(),
                    })
                status = 2
                if (query.likeType == 1) {
                    let user = await db.collection("users").where({ openId: query.openId1 }).get()
                    changeData(query, user)
                }
            } else if (operation == 0) {
                await db.collection('match').doc(matchList1.data[0]._id)
                    .update({
                        status: 3,
                        updatedAt: new Date(),
                    })
                status = 3
            }
        } else if (matchList.data.length == 0 && matchList1.data.length == 0) {
            await db.collection('match').add({
                channel: channel,
                openId1: openId1,
                openId2: openId2,
                createdAt: new Date(),
                updatedAt: new Date(),
                status: 1,
                likeType: query.likeType
            })
            status = 1
            if (query.likeType == 2) {
                let user = await db.collection("users").where({ openId: query.openId2 }).get()
                changeData(query, user)
            }
            await db.collection('information_monitor').add({
                openId: openId2,
                source: "match",
                createdAt: new Date()
            })
        }
        return res.json(ok(status))
    } catch (e) { console.log(e) }
})

//删除匹配关系
router.post("/del", async (req, res) => {
    try {
        const query = req.body
        if (query.channel) {
            let res = await db.collection('match').where({ channel: query.channel }).remove()
            return res.json(ok("解除匹配成功"))
        } else {
            return res.json(401, "参数错误")
        }
    } catch (e) { console.log(e) }
})

//获取喜欢我的人列表
router.post("/getLikeMatch", async (req, res) => {
    try {
        const query = req.body
        if (query.openId) {
            const userQuery = await db.collection('match').where({ openId2: query.openId, status: 1 }).orderBy('createdAt', 'desc').get()
            let list = []
            if (userQuery.data.length > 0) {
                for (let i = 0; i < userQuery.data.length; i++) {
                    userQuery.data[i].score = 0
                    if (userQuery.data[i].likeType == 2) {
                        userQuery.data[i].score = userQuery.data[i].score + 100
                    }
                    let person = await db.collection('users').where({ openId: userQuery.data[i].openId1 }).get()
                    let detail = await db.collection('detail_record').where({ openId: userQuery.data[i].openId1, status: 1 }).get()
                    if (detail.data.length > 0) {
                        if (detail.data[0].image.length > 0) {
                            userQuery.data[i].score = userQuery.data[i].score + 35
                        }
                    }
                    if (person.data.length == 1) {
                        if (person.data[0].image.length > 0) {
                            userQuery.data[i].score = userQuery.data[i].score + 35
                        }
                        list.push({
                            openId1: userQuery.data[i].openId1,
                            openId2: userQuery.data[i].openId2,
                            channel: userQuery.data[0].channel,
                            avatar: person.data[0].avatar,
                            name: person.data[0].name,
                            likeTime: userQuery.data[i].createdAt,
                            likeType: userQuery.data[i].likeType,
                            score: userQuery.data[i].score,
                            updatedAt: userQuery.data[i].createdAt
                        })
                    }
                }
            }
            let userList = sort(list)
            return res.json(ok(userList))
        } else {
            return res.json(fail(401, "参数错误"))
        }
    } catch (e) { console.log(e) }
})

export default router;

async function changeData(event, user) {
    if (event.likeType == 1) {
        let data = user.data[0].belike + 1
        await db.collection("users").doc(user.data[0]._id).update({
            belike: data
        })
    } else if (event.likeType == 2) {
        let data = user.data[0].besuperLike + 1
        console.log(data)
        await db.collection("users").doc(user.data[0]._id).update({
            besuperLike: data
        })

    }
}