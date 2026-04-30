# API 接口文档

## 文档概述 & 基础配置

### 基础信息
- **Base URL**: `http://localhost:8080/api`
- **协议**: HTTP/HTTPS
- **数据格式**: JSON
- **字符编码**: UTF-8

### 认证方式
- 登录/注册接口无需认证
- 其他接口通过 JWT Token 认证
- Token 有效期：24小时

---

## 统一响应格式 & 状态码

### 成功响应
```json
{
  "message": "操作成功",
  "data": { ... }
}
```

### 错误响应
```json
{
  "error": "错误信息描述"
}
```

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 参数错误/业务逻辑错误 |
| 401 | 未授权/Token无效 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 业务接口记录

### 1. 用户认证

#### 1.1 用户注册
- **URL**: `POST /user/register`
- **认证**: 无需
- **请求参数**:
  ```json
  {
    "username": "string (必填)",
    "password": "string (必填)"
  }
  ```
- **响应**:
  ```json
  {
    "message": "注册成功",
    "user": {
      "id": 1,
      "user_name": "string",
      "created_at": "2026-04-29T22:00:00Z"
    },
    "token": "jwt_token_string"
  }
  ```

#### 1.2 用户登录
- **URL**: `POST /user/login`
- **认证**: 无需
- **请求参数**:
  ```json
  {
    "username": "string (必填)",
    "password": "string (必填)"
  }
  ```
- **响应**: 同注册接口

---

### 2. 餐厅管理

#### 2.1 获取餐厅列表（支持搜索/排序）
- **URL**: `GET /restaurants`
- **认证**: 无需
- **查询参数**:
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | lat | float | 否 | 用户纬度 |
  | lon | float | 否 | 用户经度 |
  | search | string | 否 | 搜索关键词 |
  | sort | string | 否 | 排序方式（distance/score/reviews/recommended） |

- **响应**:
  ```json
  [
    {
      "id": 1,
      "name": "string",
      "latitude": 30.5185,
      "longitude": 114.3542,
      "avg_score": 4.5,
      "review_count": 128,
      "distance": 1.2,
      "final_score": 3.45
    }
  ]
  ```

#### 2.2 获取附近推荐餐厅
- **URL**: `GET /restaurants/nearby`
- **认证**: 无需
- **查询参数**:
  | 参数 | 类型 | 必填 | 说明 |
  |------|------|------|------|
  | lat | float | 是 | 用户纬度 |
  | lon | float | 是 | 用户经度 |

- **响应**: 餐厅列表（按距离排序）

#### 2.3 获取餐厅详情
- **URL**: `GET /restaurants/:id`
- **认证**: 无需
- **路径参数**: `id` - 餐厅ID
- **响应**: 单个餐厅对象

#### 2.4 获取餐厅评分列表
- **URL**: `GET /restaurants/:id/ratings`
- **认证**: 无需
- **路径参数**: `id` - 餐厅ID
- **响应**:
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "restaurant_id": 1,
      "stars": 4.5,
      "comment": "string",
      "created_at": "2026-04-29T22:00:00Z",
      "user": {
        "id": 1,
        "user_name": "string"
      }
    }
  ]
  ```

#### 2.5 创建餐厅
- **URL**: `POST /restaurants`
- **认证**: 需要
- **请求参数**:
  ```json
  {
    "name": "string (必填)",
    "latitude": 30.5185,
    "longitude": 114.3542,
    "category": "string"
  }
  ```
- **响应**: 创建的餐厅对象

---

### 3. 评分管理

#### 3.1 提交评分
- **URL**: `POST /rating`
- **认证**: 需要
- **请求参数**:
  ```json
  {
    "username": "string (必填)",
    "restaurant_id": "int (可选)",
    "restaurant_name": "string (可选)",
    "stars": "float (必填, 1-5)",
    "comment": "string (必填)"
  }
  ```
- **响应**:
  ```json
  {
    "message": "评价成功！"
  }
  ```

---

## 请求参数枚举说明

### 排序方式 (sort)
| 值 | 说明 |
|----|------|
| distance | 按距离升序 |
| score | 按评分降序 |
| reviews | 按评论数降序 |
| recommended | 智能推荐（综合评分） |

### 评分范围 (stars)
- 最小值：1.0
- 最大值：5.0
- 精度：0.1

---

## 业务特殊逻辑说明

### 1. 无定位处理
- 当请求未提供 `lat` 和 `lon` 参数时：
  - 不计算距离（distance 返回 -1）
  - 综合得分计算调整为：`评分权重(0.7) + 人气权重(0.3)`
  - 默认排序改为按评分降序

### 2. 缓存策略
- **缓存命中**：直接返回缓存数据
- **缓存未命中**：查询数据库 → 写入缓存 → 返回数据
- **缓存失效**：
  - 餐厅详情/评论：10分钟
  - 列表/推荐/搜索：2小时
  - 创建餐厅/提交评分时清除相关缓存

### 3. 餐厅识别逻辑
- 提交评分时支持通过 `restaurant_id` 或 `restaurant_name` 识别餐厅
- 优先使用 `restaurant_id`，若为 0 则使用 `restaurant_name` 查询

### 4. 评分同步机制
- 提交评分时使用数据库事务确保原子性
- 自动更新餐厅的 `avg_score` 和 `review_count`
