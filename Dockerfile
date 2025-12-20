# 构建阶段
FROM golang:1.25.5 AS builder

WORKDIR /build

# 复制 go mod 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 安装 swag CLI 工具
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制源代码
COPY . .

# 生成 Swagger 文档
RUN swag init -g ./cmd/server/main.go -o ./docs

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-bridge ./cmd/server

# 运行阶段
FROM debian:stable-slim

# 安装 CA 证书以支持 HTTPS 请求
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/notification-bridge .

# 暴露端口
EXPOSE 8000

# 运行应用
CMD ["./notification-bridge"]
