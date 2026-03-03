# 构建阶段
FROM docker.m.daocloud.io/golang:1.23-alpine as builder
WORKDIR /app

# 设置Go环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o main ./cmd/main.go

# 运行阶段
FROM docker.m.daocloud.io/alpine:latest
WORKDIR /app

# 复制构建好的二进制文件
COPY --from=builder /app/main .

RUN mkdir -p /app/config

# 复制配置文件
COPY config/config.yaml config/

# 暴露端口
EXPOSE 8000

# 启动应用
CMD ["./main"]
