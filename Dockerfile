# 第一阶段：构建阶段 - 使用官方 Go 镜像构建应用
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /book_management

# 复制 go.mod 和 go.sum 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源代码
COPY . .

# 构建应用
# -ldflags="-s -w" 用于减小二进制文件大小
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /book_management/app ./cmd

# 第二阶段：运行阶段 - 使用轻量级 Alpine 镜像运行应用
FROM alpine

# 安装必要的依赖（如你的应用需要）
# RUN apk --no-cache add ca-certificates

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /book_management/app /book_management/app

# 设置工作目录
WORKDIR /book_management

# 暴露应用运行的端口（根据你的应用调整）
EXPOSE 8989

# 设置容器启动时运行的命令
ENTRYPOINT ["./app","-conf","configs/config.yaml"]