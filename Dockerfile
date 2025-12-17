# 二开推荐阅读[如何提高项目构建效率](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/src/scene/build/speed.html)
# 建议使用 LTS 版本，如 node:20-alpine
FROM node:24-alpine

# 设置工作目录
WORKDIR /app

# 1. 设置镜像源并安装系统依赖 (合并 tzdata 和 ca-certificates)
# 容器默认时区为UTC，这里直接一并处理时区设置
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --update --no-cache \
    ca-certificates \
    tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 设置环境变量 TZ，确保 Node 应用也能识别时区
ENV TZ=Asia/Shanghai

# 2. 优先复制依赖描述文件 (利用 Docker 缓存层)
COPY package*.json ./

# 3. 安装 Node 依赖
RUN npm config set registry https://mirrors.cloud.tencent.com/npm/ \
    && npm install --production

# 4. 最后复制所有源代码
COPY . /app

# 启动命令
CMD ["npm", "start"]