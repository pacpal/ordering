# 网上订餐管理系统

## 项目概述

网上订餐管理系统是一个基于 Go 语言开发的全栈 Web 应用，提供在线点餐、购物车管理、订单处理等完整功能。系统分为前台（用户端）和后台（管理端）两个界面，支持多角色权限控制。

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 后端框架 | Gin | 高性能 Go Web 框架 |
| ORM | GORM | Go 语言 ORM 库 |
| 数据库 | SQLite | 轻量级嵌入式数据库（通过 glebarez/sqlite 纯 Go 驱动） |
| 认证 | JWT (golang-jwt/jwt/v5) | 基于 Token 的身份认证 |
| 密码加密 | bcrypt (golang.org/x/crypto) | 安全的密码哈希算法 |
| 前端 | 原生 HTML/CSS/JavaScript | 无框架依赖，单页面应用 |

## 项目结构

```
example/
├── main.go                  # 程序入口
├── config/
│   ├── app.go               # 应用配置（端口、数据库名、JWT密钥）
│   └── seed.go              # 初始数据填充
├── controllers/
│   ├── user_controller.go   # 用户注册/登录/信息
│   ├── category_controller.go # 分类 CRUD
│   ├── dish_controller.go   # 菜品 CRUD / 图片上传
│   ├── cart_controller.go   # 购物车操作
│   └── order_controller.go  # 订单创建/查询/状态管理
├── middleware/
│   └── auth.go              # JWT 认证中间件 + 管理员权限中间件
├── models/
│   └── models.go            # 数据模型定义 + 数据库连接
├── routes/
│   └── routes.go            # 路由注册
├── static/
│   └── uploads/             # 菜品图片上传目录
├── templates/
│   ├── index.html           # 前台页面（菜单浏览/购物车/下单）
│   └── admin.html           # 后台管理页面
├── go.mod
└── go.sum
```

## 快速开始

### 环境要求

- Go 1.25+

### 安装与运行

```bash
# 进入项目目录
cd example

# 安装依赖
go mod tidy

# 编译运行
go build -o ordering-system.exe .
./ordering-system.exe
```

### 访问地址

| 页面 | 地址 |
|------|------|
| 前台首页 | http://localhost:8080 |
| 后台管理 | http://localhost:8080/admin |

### 默认账号

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |

> 系统首次启动时会自动创建管理员账号、4 个菜品分类和 12 道示例菜品。

## 数据模型

### User — 用户

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| username | string(50) | 用户名，唯一索引 |
| password | string(255) | 密码（bcrypt 加密，JSON 输出时隐藏） |
| role | string(20) | 角色：`customer`（普通用户）/ `admin`（管理员） |
| phone | string(20) | 手机号 |
| address | string(255) | 地址 |

### Category — 分类

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string(50) | 分类名称，唯一索引 |
| sort | int | 排序权重（升序） |

### Dish — 菜品

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string(100) | 菜品名称 |
| price | float64 | 价格 |
| image | string(255) | 图片路径 |
| desc | string(500) | 描述 |
| status | bool | 上架状态，默认 true |
| category_id | uint | 所属分类 ID |
| category | Category | 关联分类（外键） |

### CartItem — 购物车项

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户 ID（索引） |
| dish_id | uint | 菜品 ID |
| count | int | 数量，默认 1 |
| dish | Dish | 关联菜品（外键） |

### Order — 订单

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 用户 ID（索引） |
| total | float64 | 订单总额 |
| status | string(20) | 订单状态，默认 pending |
| address | string(255) | 送餐地址 |
| phone | string(20) | 联系电话 |
| remark | string(500) | 备注 |
| items | []OrderItem | 订单明细（外键关联） |
| user | User | 关联用户（外键） |

### OrderItem — 订单明细

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| order_id | uint | 订单 ID（索引） |
| dish_id | uint | 菜品 ID |
| dish_name | string(100) | 菜品名称（下单时快照） |
| price | float64 | 单价（下单时快照） |
| count | int | 数量 |

### 订单状态流转

```
pending（待处理）→ confirmed（已确认）→ preparing（制作中）→ delivering（配送中）→ completed（已完成）
   ↓
cancelled（已取消）  ← 仅 pending 状态可由用户取消
```

## API 接口文档

### 公开接口（无需认证）

| 方法 | 路径 | 说明 | 请求参数 |
|------|------|------|----------|
| POST | /api/register | 用户注册 | `{username, password, phone?, address?}` |
| POST | /api/login | 用户登录 | `{username, password}` |
| GET | /api/dishes | 菜品列表 | `?category_id=&keyword=`（可选查询参数） |
| GET | /api/dishes/:id | 菜品详情 | — |
| GET | /api/categories | 分类列表 | — |

### 用户接口（需登录）

| 方法 | 路径 | 说明 | 请求参数 |
|------|------|------|----------|
| GET | /api/user | 获取当前用户信息 | — |
| PUT | /api/user | 更新个人信息 | `{phone, address}` |
| GET | /api/cart | 获取购物车 | — |
| POST | /api/cart | 添加到购物车 | `{dish_id, count?}` |
| PUT | /api/cart/:id | 修改购物车数量 | `{count}`（count ≥ 1） |
| DELETE | /api/cart/:id | 移除购物车项 | — |
| DELETE | /api/cart | 清空购物车 | — |
| POST | /api/orders | 创建订单 | `{address?, phone?, remark?}` |
| GET | /api/orders | 我的订单列表 | — |
| PUT | /api/orders/:id/cancel | 取消订单 | — |

### 管理员接口（需管理员权限）

| 方法 | 路径 | 说明 | 请求参数 |
|------|------|------|----------|
| GET | /api/admin/users | 用户列表 | — |
| POST | /api/admin/categories | 创建分类 | `{name, sort?}` |
| PUT | /api/admin/categories/:id | 更新分类 | `{name, sort}` |
| DELETE | /api/admin/categories/:id | 删除分类 | — |
| POST | /api/admin/dishes | 创建菜品 | `{name, price, category_id, image?, desc?}` |
| PUT | /api/admin/dishes/:id | 更新菜品 | `{name?, price?, category_id?, image?, desc?, status?}` |
| DELETE | /api/admin/dishes/:id | 删除菜品 | — |
| POST | /api/admin/upload | 上传菜品图片 | `multipart/form-data: image` |
| GET | /api/admin/orders | 所有订单列表 | — |
| PUT | /api/admin/orders/:id/status | 更新订单状态 | `{status}` |

### 认证方式

所有需认证的接口支持两种方式传递 Token：

1. **Authorization 请求头**：`Authorization: Bearer <token>`
2. **Cookie**：登录/注册成功后自动设置 `token` Cookie（有效期 24 小时）

### 通用响应格式

**成功响应**：HTTP 2xx + JSON 数据

**错误响应**：HTTP 4xx/5xx + JSON

```json
{
  "error": "错误描述"
}
```

## 功能模块

### 前台功能

- **用户注册/登录**：支持注册新用户，登录后自动设置 Token
- **菜品浏览**：按分类筛选、关键字搜索菜品
- **购物车**：添加菜品、调整数量、移除/清空
- **下单**：填写地址/电话/备注后提交订单
- **订单管理**：查看个人订单列表、取消待处理订单

### 后台功能

- **数据概览**：菜品总数、订单总数、营业额、待处理订单统计
- **分类管理**：新增、编辑、删除菜品分类
- **菜品管理**：新增、编辑、上架/下架、删除菜品，支持图片上传
- **订单管理**：查看所有订单、更新订单状态
- **用户管理**：查看注册用户列表

## 配置说明

配置文件位于 `config/app.go`：

```go
var AppConfig = Config{
    Port:      ":8080,             // 服务监听端口
    DBName:    "ordering.db",      // SQLite 数据库文件名
    JWTSecret: "your-secret-key",  // JWT 签名密钥
}
```

## 初始数据

系统首次启动时（数据库为空），`config/seed.go` 会自动填充以下数据：

- **管理员账号**：admin / admin123
- **4 个分类**：中餐、西餐、饮品、甜点
- **12 道菜品**：宫保鸡丁、红烧肉、鱼香肉丝、麻婆豆腐、牛排、意大利面、凯撒沙拉、珍珠奶茶、鲜榨橙汁、拿铁咖啡、提拉米苏、芒果布丁