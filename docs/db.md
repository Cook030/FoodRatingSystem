# 数据库设计文档

## 系统业务概述

美食评分系统是一个基于地理位置的餐厅评分推荐系统，支持用户定位、餐厅搜索、评分评论和智能推荐功能。

### 核心功能
- 用户注册/登录认证
- 餐厅信息管理（增删改查）
- 用户评分与评论
- 基于地理位置的餐厅推荐
- 多维度排序（距离/评分/评论数/智能推荐）

---

## E-R 模型设计

### 实体关系
```
User (1) ----< (N) Rating >---- (1) Restaurant
```

- **User** 与 **Rating**：一对多关系（一个用户可以发表多条评论）
- **Restaurant** 与 **Rating**：一对多关系（一个餐厅可以有多条评论）
- **Rating** 是关联表，记录用户对餐厅的评分

---

## 数据表详细结构设计

### 1. users 表（用户表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | uint (自增主键) | PRIMARY KEY | 用户唯一标识 |
| user_name | varchar(255) | UNIQUE, NOT NULL | 用户名（登录账号） |
| password_hash | varchar(255) | NOT NULL | 密码哈希值（bcrypt） |
| created_at | timestamp | AUTO | 注册时间 |

### 2. restaurants 表（餐厅表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | uint (自增主键) | PRIMARY KEY | 餐厅唯一标识 |
| name | varchar(255) | INDEX, NOT NULL | 餐厅名称 |
| latitude | decimal(10,7) | - | 纬度坐标 |
| longitude | decimal(10,7) | - | 经度坐标 |
| avg_score | decimal(3,2) | DEFAULT 0 | 平均评分 |
| category | varchar(100) | - | 餐厅分类 |
| created_at | timestamp | AUTO | 创建时间 |
| review_count | int | DEFAULT 0 | 评论数量 |

### 3. ratings 表（评分表）

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | uint (自增主键) | PRIMARY KEY | 评分唯一标识 |
| user_id | uint | INDEX, NOT NULL | 用户ID（外键） |
| restaurant_id | uint | INDEX, NOT NULL | 餐厅ID（外键） |
| stars | decimal(2,1) | - | 评分（1-5星） |
| comment | text | - | 评论内容 |
| created_at | timestamp | AUTO | 评论时间 |

---

## 表关系说明

### 外键关系
- `ratings.user_id` → `users.id`
- `ratings.restaurant_id` → `restaurants.id`

### 索引设计
| 表名 | 索引字段 | 索引类型 | 用途 |
|------|----------|----------|------|
| users | user_name | UNIQUE | 用户名唯一性校验 |
| restaurants | name | NORMAL | 餐厅名称搜索 |
| restaurants | category | NORMAL | 分类筛选 |
| ratings | user_id | NORMAL | 用户评论查询 |
| ratings | restaurant_id | NORMAL | 餐厅评论查询 |
| ratings | (user_id, restaurant_id) | COMPOSITE | 联合查询优化 |

---

## Redis 缓存设计

### 缓存策略

| 缓存键模式 | 过期时间 | 说明 |
|------------|----------|------|
| `restaurant:{id}` | 10分钟 | 餐厅详情缓存 |
| `ratings:{id}` | 10分钟 | 餐厅评论列表缓存 |
| `nearby:{lat}:{lon}` | 2小时 | 附近餐厅列表 |
| `recommend:{lat}:{lon}` | 2小时 | 推荐餐厅列表 |
| `search:{lat}:{lon}:{search}:{sort}:{hasLocation}` | 2小时 | 搜索结果缓存 |

### 缓存清除时机
- 创建新餐厅时：清除 `recommend:*`、`nearby:*`、`search:*`
- 提交评分时：清除对应餐厅的详情和评论缓存，以及所有列表缓存

---

## 业务算法与数据规则

### 智能推荐算法
```
FinalScore = 评分权重(0.6) + 距离权重(0.3) + 人气权重(0.1)
```

| 因素 | 计算方式 | 权重 |
|------|----------|------|
| 评分分数 | `AverageScore * 0.6` | 60% |
| 距离分数 | `(1/(dist+1)) * 0.3` | 30% |
| 人气分数 | `log10(ReviewCount+1) * 0.1` | 10% |

### 无定位时的排序规则
- 默认按评分降序排列
- 综合得分计算：`评分权重(0.7) + 人气权重(0.3)`

### 距离计算
使用 Haversine 公式计算两点间地球表面距离（单位：公里）

---

## 初始化 SQL 脚本

### 测试用户
```sql
INSERT INTO users (user_name, password_hash, created_at) VALUES
('u1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', NOW());
```

### 示例餐厅数据
详见 `server/database/insert_wuhan_data.sql`，包含武汉地区 30+ 家特色餐厅数据。
