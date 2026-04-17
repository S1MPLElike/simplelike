### 项目简介
这是一个前后端分离的社交平台系统，支持用户注册登录、关注互动、帖子发布、实时聊天等功能。后端采用Go语言开发，前端使用Vue3构建。

### 技术栈

**后端**
- **Go 1.21+** - 高性能服务端语言
- **Gin框架** - 轻量级Web框架，高并发支持
- **GORM** - Go语言ORM库，简化数据库操作
- **MySQL 8.0** - 关系型数据库
- **WebSocket** - 实时双向通信
- **Token + Cookie** - 用户身份认证

**前端**
- Vue 3 + Composition API
- Vite 构建工具
- Axios HTTP 客户端

### 后端架构设计

#### 1. 分层架构
```
router/      → 路由层，统一注册API路由
controller/  → 控制器层，处理业务逻辑
model/       → 数据模型层，定义数据结构
utils/       → 工具函数层，封装通用功能
config/      → 配置管理层，集中管理配置
```

#### 2. 核心模块

**用户认证模块**
- Token生成与验证机制
- Token生成后存到客户端的cookie中
- 登录状态7前到期，到期前用户点开该web会直接验证Token自动登录

**关注关系模块**
- 采用`follows`表维护关注关系
- 支持关注/取消关注操作
- 实现互相关注自动成为好友的业务逻辑

**帖子模块**
- 完整的CRD操作（没有U）
- 点赞、收藏、评论功能
- 支持搜索和热门帖子排序
- 实现用户主页帖子分页查询

**实时聊天模块**
- 基于WebSocket的实时通信
- 支持文本、图片、emoji消息类型
- 消息持久化存储（MySQL）
- 未读消息计数

### API设计

采用RESTful风格设计，统一的响应格式：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

主要接口分组：
- `/api/v1/user/*` - 用户相关
- `/api/v1/post/*` - 帖子相关  
- `/api/v1/chat/*` - 聊天相关

### 数据库设计

**核心表结构**
- `users` - 用户表
- `posts` - 帖子表
- `comments` - 评论表
- `likes` - 点赞表
- `collects` - 收藏表
- `follows` - 关注关系表
- `messages` - 聊天消息表

使用GORM的AutoMigrate自动迁移表结构，支持字段自动更新。

### 项目亮点（后端）

1. **高性能**：Go语言天然支持高并发，Gin框架处理速度快
2。 **ORM封装**：使用GORM简化数据库操作，支持链式调用
3。 **实时通信**：WebSocket实现消息实时推送，支持双向通信
4。 **数据一致性**：部分对实时性要求高的场景直接访问Mysql，部分对性能要求高的场景用redis作为缓存层来加快速度
5。 **业务完整**：涵盖社交平台核心功能，用户、帖子、关系、消息

### 部署方式

后端：
```bash
cd Back
go mod tidy
go install github.com/cosmtrek/air@latest
air
```

前端：
```bash
cd Front
npm install
npm run dev
```

---
