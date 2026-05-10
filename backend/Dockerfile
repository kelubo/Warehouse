# 多阶段构建
FROM golang:1.22-alpine AS builder

WORKDIR /app/backend

# 安装构建依赖
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

# 构建应用
RUN go build -o warehouse-management ./cmd/main.go

# 生产镜像
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/backend/warehouse-management .
COPY backend/.env .

EXPOSE 8080

CMD ["./warehouse-management"]