# 餐厅评分系统 (Food Rating System)

一个基于地理位置的餐厅评分推荐系统，支持用户定位、餐厅搜索、评分评论和智能推荐功能。

## 目录结构

```
foodRatingSystem/
├── client/                    # 前端 (React + TypeScript)
│   ├── src/
│   │   ├── api/              # API 请求封装
│   │   │   └── index.ts      # 餐厅、评分 API 方法
│   │   ├── components/        # 公共组件
│   │   │   ├── MapPanel.tsx      # 地图面板组件
│   │   │   ├── RatingForm.tsx    # 评分表单组件
│   │   │   ├── RestaurantCard.tsx  # 餐厅卡片组件
│   │   │   └── RestaurantForm.tsx  # 新增餐厅表单
│   │   ├── hooks/            # 自定义 Hooks
│   │   │   └── useLocation.ts    # 地理位置定位 Hook
│   │   ├── pages/            # 页面组件
│   │   │   ├── Detail.tsx       # 餐厅详情页
│   │   │   └── Home.tsx          # 首页
│   │   ├── types/            # TypeScript 类型定义
│   │   ├── App.tsx           # 根组件
│   │   ├── main.tsx          # 入口文件
│   │   └── index.css         # 全局样式
│   ├── package.json
│   └── tailwind.config.js    # Tailwind CSS 配置
│
└── server/                    # 后端 (Go + Gin + GORM)
    ├── config/               # 配置管理
    │   └── config.go         # 配置文件加载
    ├── database/             # 数据库连接
    │   └── db.go             # GORM 数据库连接
    ├── handler/              # HTTP 请求处理
    │   ├── rating.go         # 评分相关 Handler
    │   └── restaurant.go     # 餐厅相关 Handler
    ├── middleware/           # 中间件
    │   └── cors.go           # CORS 跨域配置
    ├── model/                # 数据模型
    │   ├── rating.go         # 评分数据模型
    │   └── restaurant.go     # 餐厅数据模型
    ├── repository/           # 数据访问层
    │   ├── rating_repo.go    # 评分数据操作
    │   └── restaurant_repo.go  # 餐厅数据操作
    ├── router/               # 路由配置
    │   └── router.go         # Gin 路由设置
    ├── service/              # 业务逻辑层
    │   ├── rating_service.go    # 评分业务逻辑
    │   └── restaurant_service.go # 餐厅业务逻辑
    ├── utils/                # 工具函数
    │   └── distance.go       # 距离计算工具
    ├── main.go              # 入口文件
    └── go.mod
```

---

## 开发任务索引

### 已完成功能

- [x] **用户地理位置定位** - 获取用户当前经纬度
- [x] **餐厅列表展示** - 首页展示附近餐厅卡片
- [x] **餐厅搜索** - 按名称搜索餐厅（带 3000ms 防抖）
- [x] **多维度排序** - 按距离/评分/评论数/综合推荐排序
- [x] **餐厅详情页** - 查看餐厅详细信息及评分列表
- [x] **提交评分** - 用户对餐厅进行 1-5 星评分和评论
- [x] **新增餐厅** - 通过地图选点创建新餐厅
- [x] **智能推荐算法** - 基于评分、距离、人气权重的综合推荐
- [x] **地图交互** - 显示餐厅位置、点击地图选点

---

## 核心技术实现

### 1. 综合推荐算法

**位置**: `server/service/restaurant_service.go`

**思路**: 基于加权评分的推荐算法，综合考虑三个因素：

```
FinalScore = 评分权重(0.6) + 距离权重(0.3) + 人气权重(0.1)
```

| 因素 | 计算方式 | 权重 |
|------|----------|------|
| 评分分数 | `AverageScore * 0.6` | 60% |
| 距离分数 | `(1/(dist+1)) * 0.3` | 30% |
| 人气分数 | `log10(ReviewCount+1) * 0.1` | 10% |

**设计理由**:
- 评分权重最高（民以食为天，口碑最重要）
- 距离使用 `1/(dist+1)` 确保持续正值且递减
- 人气使用 `log10` 平滑处理，防止大店霸榜

### 2. 距离计算

**位置**: `server/utils/distance.go`

**思路**: 使用 Haversine 公式计算两点间地球表面距离

```
a = sin²(Δlat/2) + cos(lat1) * cos(lat2) * sin²(Δlng/2)
c = 2 * atan2(√a, √(1-a))
d = R * c  (R = 6371km)
```

### 3. 数据库评分同步

**位置**: `server/repository/rating_repo.go`

**思路**: 事务确保评分插入和餐厅平均分更新原子性

```go
func AddRatingAndUpdateScore(r model.Rating) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 插入新评分
        if err := tx.Create(&r).Error; err != nil {
            return err
        }
        // 2. 计算新的平均分和评论数
        // 3. 更新餐厅记录
        return tx.Model(&model.Restaurant{}).Where("id = ?", r.RestaurantID).Updates(...).Error
    })
}
```

### 4. 前端搜索防抖

**位置**: `client/src/pages/Home.tsx`

**思路**: 使用 `setTimeout` 延迟搜索请求，避免频繁 API 调用

```typescript
useEffect(() => {
    const timeout = setTimeout(() => {
        fetchRestaurants(search, sortBy);
    }, 300);
    return () => clearTimeout(timeout);
}, [search, sortBy]);
```

### 5. RESTful API 设计

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/restaurants` | 获取餐厅列表（支持搜索、排序） |
| GET | `/api/restaurants/nearby` | 获取附近推荐餐厅 |
| GET | `/api/restaurants/:id` | 获取餐厅详情 |
| GET | `/api/restaurants/:id/ratings` | 获取餐厅评分列表 |
| POST | `/api/restaurants` | 创建新餐厅 |
| POST | `/api/rating` | 提交评分 |

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
| 架构 | Repository-Service-Handler 三层架构 |
