# 仓库管理系统

基于 Go + Gin + Flutter 的仓库管理系统，支持多数据库（MySQL/PostgreSQL/SQLite），配套移动端应用。

## 功能特性

- 用户管理（注册、登录、JWT认证）
- 产品管理（增删改查）
- 仓库管理（增删改查）
- 货架管理（多层级货架）
- 箱子管理（货架内箱子）
- 库存管理（入库、出库）
- 订单管理（入库单、出库单）
- 移动端应用（支持离线使用）
- 数据同步（服务器与移动端实时同步）
- 多数据库支持（MySQL、PostgreSQL、SQLite）

## 技术栈

### 后端
- Go 1.25+
- Gin 框架
- GORM ORM
- MySQL / PostgreSQL / SQLite
- JWT 认证
- Swagger API 文档

### 移动端
- Flutter
- SQLite 本地存储
- HTTP 通信
- 数据同步机制

## 项目结构

```
D:\Git\Warehouse\
├── backend/              # Go 后端代码
│   ├── cmd/             # 入口命令
│   │   ├── main.go      # 主程序入口
│   │   └── initialize/  # 数据库初始化
│   ├── config/          # 配置管理
│   ├── controllers/     # 控制器
│   ├── database/        # 数据库配置
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── routes/          # 路由配置
│   ├── server/          # 服务器
│   ├── utils/           # 工具函数
│   ├── static/          # 静态文件
│   ├── .env             # 环境配置
│   ├── Dockerfile       # Docker配置
│   └── docker-compose.yml
│
├── mobile/               # Flutter 移动端
│   ├── lib/
│   │   ├── models/      # 数据模型
│   │   ├── screens/     # 页面
│   │   ├── services/    # 服务层
│   │   └── main.dart
│   └── pubspec.yaml
│
├── LICENSE
└── README.md
```

## 快速开始

### 环境要求

- Go 1.25+
- Flutter SDK（移动端开发）
- MySQL 5.7+ / PostgreSQL 12+ / SQLite

### 后端启动

```bash
# 进入后端目录
cd backend

# 下载依赖
go mod download

# 初始化数据库（首次运行）
go run ./cmd/main.go init

# 启动服务
go run ./cmd/main.go
```

服务将在 http://localhost:8080 启动。

API 文档：http://localhost:8080/swagger/index.html

### 移动端启动

```bash
# 进入移动端目录
cd mobile

# 下载依赖
flutter pub get

# 运行应用
flutter run
```

### 默认管理员

- 用户名：admin
- 密码：admin123

## 配置说明

### 环境变量 (.env)

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| ENV | development | 运行环境 |
| PORT | 8080 | 服务端口 |
| DB_TYPE | mysql | 数据库类型 |
| DB_HOST | localhost | 数据库主机 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USER | root | 数据库用户名 |
| DB_PASSWORD | - | 数据库密码 |
| DB_NAME | warehouse | 数据库名称 |
| DATABASE_PATH | ./warehouse.db | SQLite数据库路径 |
| JWT_SECRET | warehouse-secret-key | JWT密钥 |
| JWT_EXPIRE_HOURS | 24 | JWT过期时间 |
| LOG_LEVEL | info | 日志级别 |

### 支持的数据库类型

- **MySQL**: 设置 `DB_TYPE=mysql`
- **PostgreSQL**: 设置 `DB_TYPE=postgres`
- **SQLite**: 设置 `DB_TYPE=sqlite`

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

### 货架管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/shelves | 获取货架列表 |
| GET | /api/v1/shelves/:id | 获取货架详情 |
| POST | /api/v1/shelves | 创建货架 |
| PUT | /api/v1/shelves/:id | 更新货架 |
| DELETE | /api/v1/shelves/:id | 删除货架 |

### 箱子管理

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/v1/boxes | 获取箱子列表 |
| GET | /api/v1/boxes/:id | 获取箱子详情 |
| POST | /api/v1/boxes | 创建箱子 |
| PUT | /api/v1/boxes/:id | 更新箱子 |
| DELETE | /api/v1/boxes/:id | 删除箱子 |

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

## Docker 部署

```bash
# 进入后端目录
cd backend

# 构建并启动
docker-compose up -d
```

## 开发命令

- `go run ./cmd/main.go` - 启动服务
- `go run ./cmd/main.go init` - 初始化数据库
- `go run ./cmd/main.go help` - 显示帮助
- `go build ./cmd/main.go` - 编译项目

## License

MIT
