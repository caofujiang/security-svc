# 编译阶段
FROM golang:alpine AS build-env
LABEL maintainer="FuJiang Cao"
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /src
COPY . .
RUN go mod download && go build -o security-svc .

# 创建一个新的阶段用于运行时环境
FROM alpine:latest
# 保持时区设置
RUN apk update && apk upgrade
RUN apk add bash
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app
COPY --from=build-env /src/security-svc .
COPY --from=build-env /src/conf  ./conf/
COPY --from=build-env /src/fscan ./fscan/
RUN chmod +x   ./fscan/*

# 为运行时赋予必要的权限
RUN chmod +x ./security-svc
EXPOSE 8000
CMD ["/app/security-svc"]

