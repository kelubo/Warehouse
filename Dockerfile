# 多阶段构建
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 构建应用
RUN go build -o warehouse-management .

# 生产镜像
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/warehouse-management .
COPY .env .

EXPOSE 8080

CMD ["./warehouse-management"]