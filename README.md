# 仓库管理系统

基于 Go + Gin + SQLite 的仓库管理系统。

## 功能特性

- 用户管理（注册、登录）
- 产品管理（增删改查）
- 仓库管理（增删改查）
- 库存管理（入库、出库）
- 订单管理（入库单、出库单）
- JWT 认证

## 技术栈

- Go 1.25+
- Gin 框架
- GORM ORM
- SQLite 数据库
- JWT 认证

## 快速开始

### 环境要求

- Go 1.25+

### 安装依赖

```bash
go mod download
```

### 运行项目

```bash
go run main.go
```

服务将在 http://localhost:8080 启动。

### 默认管理员

- 用户名：admin
- 密码：admin123

## API 接口

### 认证

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/register | 用户注册 |
| POST | /api/v1/login | 用户登录 |

### 产品管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/products | 获取产品列表 |
| GET | /api/v1/products/:id | 获取产品详情 |
| POST | /api/v1/products | 创建产品 |
| PUT | /api/v1/products/:id | 更新产品 |
| DELETE | /api/v1/products/:id | 删除产品 |

### 仓库管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/warehouses | 获取仓库列表 |
| GET | /api/v1/warehouses/:id | 获取仓库详情 |
| POST | /api/v1/warehouses | 创建仓库 |
| PUT | /api/v1/warehouses/:id | 更新仓库 |
| DELETE | /api/v1/warehouses/:id | 删除仓库 |

### 库存管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/inventories | 获取库存列表 |
| GET | /api/v1/inventories/:id | 获取库存详情 |
| POST | /api/v1/inventories | 创建库存记录 |
| PUT | /api/v1/inventories/:id | 更新库存记录 |
| DELETE | /api/v1/inventories/:id | 删除库存记录 |

### 订单管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/orders | 获取订单列表 |
| GET | /api/v1/orders/:id | 获取订单详情 |
| POST | /api/v1/orders | 创建订单 |
| PUT | /api/v1/orders/:id | 更新订单 |
| DELETE | /api/v1/orders/:id | 删除订单 |

## 项目结构

```
warehouse-management/
├── controllers/     # 控制器
├── middleware/      # 中间件
├── models/          # 数据模型
├── database/        # 数据库配置
├── config/          # 配置管理
├── routes/          # 路由配置
├── main.go          # 入口文件
└── go.mod           # 依赖管理
```

## 配置说明

可以通过环境变量配置：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| ENV | development | 运行环境 |
| PORT | 8080 | 服务端口 |
| DATABASE_PATH | ./warehouse.db | 数据库路径 |
| JWT_SECRET | warehouse-secret-key | JWT密钥 |

## License

MIT
