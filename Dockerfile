# 二开推荐阅读[如何提高项目构建效率](https://developers.weixin.qq.com/miniprogram/dev/wxcloudrun/src/scene/build/speed.html)
FROM node:24-alpine
WORKDIR /app
COPY . /app
# 容器默认时区为UTC，如需使用上海时间请启用以下时区设置命令
# RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone

# 安装依赖包，如需其他依赖包，请到alpine依赖包管理(https://pkgs.alpinelinux.org/packages?name=php8*imagick*&branch=v3.13)查找。
# 选用国内镜像源以提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
&& apk add --update --no-cache nodejs npm ca-certificates\
&& npm config set registry https://mirrors.cloud.tencent.com/npm/ \
&& npm install

# --- 设置时区开始 ---
# 1. 安装 tzdata 包
RUN apk add --no-cache tzdata

# 2. 设置环境变量
ENV TZ=Asia/Shanghai

# 3. 复制时区文件 (Alpine 的做法)
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
# --- 设置时区结束 ---

CMD ["npm", "start"]