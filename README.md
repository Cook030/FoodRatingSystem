# 餐厅评分系统 (Food Rating System)

一个基于地理位置的餐厅评分推荐系统，支持用户定位、餐厅搜索、评分评论和智能推荐功能。

---

## 目录结构

```
foodRatingSystem/
├── client/                    # 前端 (React + TypeScript)
│   ├── src/
│   │   ├── api/              # API 请求封装
│   │   ├── components/        # 公共组件
│   │   ├── hooks/            # 自定义 Hooks
│   │   ├── pages/            # 页面组件
│   │   ├── types/            # TypeScript 类型定义
│   │   ├── App.tsx           # 根组件
│   │   ├── main.tsx          # 入口文件
│   │   └── index.css         # 全局样式
│   └── package.json
│
└── server/                    # 后端 (Go + Gin + GORM)
    ├── config/               # 配置管理
    ├── database/             # 数据库连接与初始化脚本
    ├── handler/              # HTTP 请求处理
    ├── middleware/           # 中间件
    ├── model/                # 数据模型
    ├── repository/           # 数据访问层
    ├── router/               # 路由配置
    ├── service/              # 业务逻辑层
    ├── utils/                # 工具函数
    └── main.go              # 入口文件
```

---

## 功能特性

- [x] **用户认证** - 注册/登录功能
- [x] **地理位置定位** - 获取用户当前经纬度
- [x] **餐厅列表展示** - 首页展示附近餐厅卡片
- [x] **餐厅搜索** - 按名称搜索餐厅（带 3000ms 防抖）
- [x] **多维度排序** - 按距离/评分/评论数/综合推荐排序
- [x] **餐厅详情页** - 查看餐厅详细信息及评分列表
- [x] **提交评分** - 用户对餐厅进行 1-5 星评分和评论
- [x] **新增餐厅** - 通过地图选点创建新餐厅
- [x] **智能推荐算法** - 基于评分、距离、人气权重的综合推荐
- [x] **地图交互** - 显示餐厅位置、点击地图选点
- [x] **Redis 缓存** - 热点数据缓存加速

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端框架 | React 18 + TypeScript |
| UI 样式 | Tailwind CSS |
| 地图 | Leaflet + React-Leaflet |
| 后端框架 | Go + Gin |
| ORM | GORM |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 架构 | Repository-Service-Handler 三层架构 |

---

## 从零开始启动项目

### 环境要求

- **Node.js** >= 16.x（前端）
- **Go** >= 1.18（后端）
- **PostgreSQL** >= 13（数据库）
- **Redis** >= 6.x（缓存）

---

### 1. 数据库准备

#### 1.1 安装并启动 PostgreSQL

**Windows**: 下载安装包 https://www.postgresql.org/download/windows/

**macOS**:
```bash
brew install postgresql
brew services start postgresql
```

**Ubuntu/Debian**:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

#### 1.2 创建数据库

```bash
psql -U postgres
```

在 psql 命令行中执行：
```sql
CREATE DATABASE food_rating_db;
\q
```

#### 1.3 修改数据库连接配置

编辑 `server/database/db.go`，修改 DSN 连接字符串：

```go
dsn := "host=localhost user=postgres password=你的密码 dbname=food_rating_db port=5432 sslmode=disable"
```

---

### 2. Redis 准备

#### 2.1 安装并启动 Redis

**Windows**: 使用 WSL 或下载 https://github.com/tporadowski/redis/releases

**macOS**:
```bash
brew install redis
brew services start redis
```

**Ubuntu/Debian**:
```bash
sudo apt install redis-server
sudo systemctl start redis
```

#### 2.2 验证 Redis 连接

```bash
redis-cli ping
# 应返回 PONG
```

---

### 3. 初始化数据

#### 3.1 启动后端（自动建表）

```bash
cd server
go run main.go
```

首次启动会自动创建数据表结构。

#### 3.2 插入示例数据

在 PostgreSQL 中执行初始化脚本：

```bash
psql -h localhost -U postgres -d food_rating_db -f server/database/insert_wuhan_data.sql
```

---

### 4. 启动后端服务

```bash
cd server
go run main.go
```

后端服务默认运行在 `http://localhost:8080`

---

### 5. 启动前端

#### 5.1 安装依赖

```bash
cd client
npm install
```

#### 5.2 启动开发服务器

```bash
npm run dev
```

前端默认运行在 `http://localhost:5173`

---

### 6. 访问项目

打开浏览器访问 `http://localhost:5173`

**测试账号**:
- 用户名: `u1`
- 密码: `123456`

---

## 项目文档

- [数据库设计文档](docs/db.md) - 数据表结构、E-R 模型、Redis 缓存设计
- [API 接口文档](docs/api.md) - 接口列表、参数说明、业务逻辑

---

## 核心算法

### 智能推荐算法

```
FinalScore = 评分权重(0.6) + 距离权重(0.3) + 人气权重(0.1)
```

| 因素 | 计算方式 | 权重 |
|------|----------|------|
| 评分分数 | `AverageScore * 0.6` | 60% |
| 距离分数 | `(1/(dist+1)) * 0.3` | 30% |
| 人气分数 | `log10(ReviewCount+1) * 0.1` | 10% |

### 距离计算

使用 Haversine 公式计算两点间地球表面距离（单位：公里）
