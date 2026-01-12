import tcb from "@cloudbase/node-sdk";//用来连接云开发的数据库

const tencent_cloud = tcb.init({ //初始化环境，后面可能直接从容器环境里面读取
    secretId: "AKIDzveAVfdOwHOHxPl5KNF1oTwELG3e4GMX",
    secretKey: "BB50AE6lFZhcKrFRQErHTREpDRUM1Vn2",
    env: "cloud1-1gth9cum37c9015c"
})

export function getDb(){
    return tencent_cloud.database()
}

export function getModels(){
    return tencent_cloud.models
}

export function getTcb(){
    return tencent_cloud
}