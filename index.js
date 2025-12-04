import express from "express";//用express框架简单一点，以后可能换成Rust或者什么的
import tcb from "@cloudbase/node-sdk";//用来连接云开发的数据库
import crypto from "crypto"; //用Nodejs的加密包
import morgan from "morgan";
import { error, trace } from "console";

const tencent_cloud = tcb.init({ //初始化环境，后面可能直接从容器环境里面读取
    secretId: "AKIDzveAVfdOwHOHxPl5KNF1oTwELG3e4GMX",
    secretKey: "BB50AE6lFZhcKrFRQErHTREpDRUM1Vn2",
    env: "cloud1-1gth9cum37c9015c"
})
const db = tencent_cloud.database();//初始化数据库，之后可能直接从模型倒下去
const models = tencent_cloud.models; //简化模型调用
const SECRET_KEY = "37b6c7156eaa44d0"; //GoEasy API请求检验
const app = express(); //创建 express instance
app.use(express.urlencoded({ extended: false })); //启用 URLEncode 反序列化
app.use(express.json()); //启用JSON反序列化
app.use(morgan('combined'));//启用morgan日志记录器
app.get("/", async (req, res) => {
    res.send("Hello");
})

// 微信云托管服务，健康检查和内容安全回调
app.post("/censor", async (req, res) => {
    let { action, MsgType, Event, trace_id, result } = req.body; //反序列化几个关键的对象

    // 1. 微信服务健康检查
    if (action === "CheckContainerPath") {
        console.log("Received wechat health check.");
        return res.status(200).send();
    }

    // 2. 微信内容安全检查回调
    if (MsgType === "event" && Event === "wxa_media_check") {
        res.status(200).send();
        models.picture_list.update({
            data: {
                secCheckStatus: result.suggest == "pass" ? 1 : 0
            },
            filter: {
                where: {
                    traceID: {
                        $eq: trace_id
                    }
                }
            }
        }).then(data => { return data })
    } else {
        // 默认回复，避免微信重试
        console.log("Received unexpected POST to /censor:", req.body);
        return res.status(200).send();
    }
});
//GoEasy回调信息
app.post("/webhook", verifySignature, async (req, res) => {
    res.json({
        "code": 200,
        "content": "success"
    });
    let content = JSON.parse(req.body.content);
    if (Array.isArray(content) && content.length > 0) {
        content.forEach(
            (message, index) => {
                let _message = JSON.parse(message.content);
                models.new_chat_history.create({
                    data: {
                        senderOpenID: _message.openID,
                        channelId: message.channel,
                        timestamp: message.timestamp,
                        messageContent: {
                            "id": _message.id,
                            "content": _message.content,
                            "type": _message.type,
                            "contentType": _message.contentType,
                            "pic": _message.pic,
                            "name": _message.name,
                            "state": _message.state
                        }
                    }
                })
            }
        )
    }
})
// 获取临时文件URL
app.post("/startCensor", async (req, res) => {
    try {
        var openid = req.headers["x-wx-openid"];
        console.log(openid);
        var { fileid, digest } = req.body;
        console.log(req.body);

        if (!fileid) {
            console.log("The request contains no fileid");
            return res.status(400).json({
                code: 400,
                message: "fileid is required"
            });
        } else {
            let uploadedFile = await models.picture_list.get({
                filter: {
                    where: {
                        picHash: {
                            $eq: digest,
                        }
                    }
                }
            });
            if (Object.entries(uploadedFile.data) == 0) {
                res.status(200).json({
                    code: 200,
                    message: "there is no same file uploaded",
                    secCheckStatus: 2,
                });
            } else {
                return res.status(200).json({
                    code: 200,
                    message: "the same file has been uploaded already",
                    secCheckStatus: uploadedFile.data.secCheckStatus
                })
            }
        }
        var fileurl_res = await tencent_cloud.getTempFileURL({
            fileList: [fileid]
        });
    } catch (error) {
        console.error("Error getting temp file URL:", error);
        res.status(500).json({
            code: 500,
            message: "Internal server error",
            error: error.message
        });
    }
    //console.log(fileurl_res);
    let wx_backend_response = await fetch("http://api.weixin.qq.com/wxa/media_check_async", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "media_url": fileurl_res.fileList[0].tempFileURL, //临时文件URL
            "media_type": 2,
            "openid": openid,
            "version": 2, //小程序版本号，默认为2
            "scene": 1 //场景值，用户资料
        })
    });
    if (!wx_backend_response.ok) {
        let error_message = wx_backend_response.json().then(data=>{console.log(data);return data});
    }
    let wx_backend_response_json = await wx_backend_response.json();
    console.log(wx_backend_response_json);
    models.picture_list.create({
        data: {
            userPicUrl: fileid,
            secCheckStatus: 2,
            traceID: wx_backend_response_json.trace_id,
            picHash: digest,
        }
    })

});

const port = process.env.PORT || 80;

async function bootstrap() {
    app.listen(port, () => {
        console.log("启动成功", port);
    });
}
async function verifySignature(req, res, next) {
    //console.log(req.body)
    const receivedSignature = req.get('x-goeasy-signature'); // 假设签名通过 header 传递
    if (!receivedSignature) {
        return res.status(401).json({ error: 'Missing signature' });
    }
    let body = req.body; //在启用反序列化之后就可以这样写了
    let content = body.content;
    console.log(body); //给出日志看一下
    // 计算预期签名：HMAC-SHA1 + Base64 编码
    let expectedSignature = crypto
        .createHmac('sha1', SECRET_KEY)
        .update(content)
        .digest('base64'); //自动使用 raw=true 并 base64 编码

    // 安全比较（防时序攻击）
    if (crypto.timingSafeEqual(
        Buffer.from(receivedSignature),
        Buffer.from(expectedSignature)
    )) {
        next(); // 校验通过
    } else {
        res.status(401).json({ error: 'Invalid signature' });
    }
}
bootstrap();
