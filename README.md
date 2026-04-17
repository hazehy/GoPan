# GoPan

GoPan 是一个自建网盘系统，支持用户注册登录、文件上传与下载、分片上传、文件管理、分享链接、管理员后台以及审计日志。项目采用前后端分离架构，后端基于 Go + go-zero，前端基于 Vue 3 + TypeScript + Vite。

## 功能概览

- 用户注册、登录、令牌刷新、用户详情查询
- 发送邮箱验证码
- 文件上传、预上传、分片上传、分片上传完成、下载、删除、重命名、移动
- 文件夹创建和用户目录树管理
- 资源保存和分享链接创建
- 管理员控制台：系统概览、用户列表、用户状态管理、文件列表、文件删除、日志查看
- 腾讯云 COS 文件存储和分片上传支持

## 技术栈

- 后端：Go 1.25.7、go-zero、xorm、MySQL、Redis、JWT
- 对象存储：腾讯云 COS
- 前端：Vue 3、TypeScript、Vite、Pinia、Vue Router、Axios
- 部署：Docker、Docker Compose、Nginx

## 代码结构

### 后端

- `gopan/gopan.go`：后端服务入口
- `gopan/gopan.api`：go-zero API 定义
- `gopan/etc/`：不同环境的配置文件
- `gopan/internal/handler/`：HTTP 路由和请求处理
- `gopan/internal/logic/`：业务逻辑
- `gopan/internal/middleware/`：认证、管理员鉴权等中间件
- `gopan/internal/svc/`：服务上下文和依赖注入
- `gopan/internal/types/`：请求与响应结构
- `gopan/helper/`：通用工具，例如 JWT、验证码、COS、密码处理
- `gopan/models/`：数据库模型

### 前端

- `web/src/views/`：页面视图，例如登录、注册、网盘、后台、分享页
- `web/src/router/`：路由和路由守卫
- `web/src/stores/`：Pinia 状态管理
- `web/src/api/`：接口封装
- `web/src/components/`：通用组件
- `web/src/utils/`：工具方法
- `web/nginx/default.conf`：生产环境 Nginx 配置

### 其他关键文件

- `databases.sql`：数据库建表脚本
- `docker-compose.yml`：完整服务编排
- `deploy-docker.sh`：Linux 部署脚本
- `.env.example`：Docker 部署环境变量示例

## 业务架构

### 文件存储模式

系统使用“公共文件池 + 用户仓库”的设计：

- `repository_pool` 保存文件内容本体，以哈希去重
- `user_repository` 保存用户自己的目录结构和文件引用
- `share_link` 保存分享信息和访问控制

这样做的好处是：

- 相同内容的文件只保留一份物理数据
- 用户可以拥有自己的目录和文件名，而不重复保存内容
- 分享和文件管理可以围绕统一的文件池展开

### 权限模型

`user_basic` 里有几类关键字段：

- `status`：用户是否正常或禁用
- `role`：普通用户或管理员
- `upload_permission`：上传权限
- `download_permission`：下载权限
- `share_permission`：分享权限

管理员接口需要额外通过 Admin 中间件校验，普通用户登录后仅能访问授权路由。

### 认证流程

- 前端登录后会把 `token` 和 `refresh_token` 存到 `localStorage`
- 请求通过 Axios 拦截器自动带上 `Authorization`
- 后端认证中间件会校验 JWT，并检查用户状态
- 前端在收到 `401` 时会自动尝试用 refresh token 刷新

### 前端路由

- `/login`：登录页
- `/register`：注册页
- `/disk`：网盘首页，需要登录
- `/admin`：后台页，需要管理员权限
- `/share/:identity`：分享页，无需登录

## 数据库表说明

`databases.sql` 中的主要表如下。

| 表名              | 作用                                     |
| ----------------- | ---------------------------------------- |
| `user_basic`      | 用户基础信息、角色、权限和登录状态       |
| `repository_pool` | 公共文件池，保存文件内容、哈希和存储路径 |
| `user_repository` | 用户目录树和文件引用                     |
| `share_link`      | 分享记录、过期时间、访问次数和密码哈希   |
| `audit_log`       | 操作审计日志                             |

### 关键关系

- `user_repository.user_identity` 关联 `user_basic.identity`
- `share_link.user_identity` 关联 `user_basic.identity`
- `user_repository.repository_identity` 指向 `repository_pool.identity`
- 文件夹节点在 `user_repository` 中表现为 `repository_identity` 为空

## 环境要求

### 本地开发

- Go 1.25.7 或兼容版本
- Node.js 20 或兼容版本
- MySQL 8.x
- Redis 7.x

### Docker 部署

- Docker Engine
- Docker Compose v2 优先，v1 也可用但不推荐

## 配置说明

### 后端配置文件

后端配置位于：

- `gopan/etc/gopan-api.yaml`
- `gopan/etc/gopan-api.local.yaml`
- `gopan/etc/gopan-api.docker.yaml`

本地开发建议使用 `gopan/etc/gopan-api.local.yaml`。Docker 部署时，`gopan/entrypoint.sh` 会把 `gopan-api.docker.yaml` 渲染成实际运行配置。

### Docker 环境变量

项目根目录提供了 `.env.example`，首次部署建议复制为 `.env` 并修改敏感信息。

| 变量                       | 说明                          |
| -------------------------- | ----------------------------- |
| `MYSQL_ROOT_PASSWORD`      | MySQL root 密码               |
| `MYSQL_DATABASE`           | 数据库名，默认 `gopan`        |
| `GOPAN_REDIS_ADDR`         | Redis 地址，默认 `redis:6379` |
| `MYSQL_IMAGE`              | MySQL 镜像版本                |
| `REDIS_IMAGE`              | Redis 镜像版本                |
| `GOPAN_JWT_KEY`            | JWT 签名密钥                  |
| `GOPAN_FROM_MAIL`          | 发件人邮箱                    |
| `GOPAN_MAIL_PASSWORD`      | 邮箱授权码                    |
| `GOPAN_SMTP_HOST`          | SMTP 主机                     |
| `GOPAN_SMTP_PORT`          | SMTP 端口                     |
| `GOPAN_TENCENT_SECRET_ID`  | 腾讯云 SecretId               |
| `GOPAN_TENCENT_SECRET_KEY` | 腾讯云 SecretKey              |
| `GOPAN_COS_BUCKET_URL`     | 腾讯云 COS Bucket URL         |
| `TZ`                       | 时区，默认 `Asia/Shanghai`    |
| `FRONTEND_PORT`            | 前端对外映射端口，默认 `80`   |

### 典型本地配置

如果你先在本地验证后端，至少要确保：

- MySQL 连接串可用
- Redis 地址正确
- `GOPAN_JWT_KEY` 已设置
- 邮件验证码功能所需的 SMTP 配置已设置
- 如果要启用 COS 上传，`GOPAN_TENCENT_SECRET_ID`、`GOPAN_TENCENT_SECRET_KEY` 和 `GOPAN_COS_BUCKET_URL` 都要正确配置

## 本地开发

### 1. 准备依赖

先启动 MySQL 和 Redis，并执行 `databases.sql` 初始化表结构。

### 2. 启动后端

在仓库根目录执行：

```bash
go run ./gopan/gopan.go -f gopan/etc/gopan-api.local.yaml
```

后端默认监听 `8888` 端口。

### 3. 启动前端

在 `web/` 目录下执行：

```bash
npm install
npm run dev
```

Vite 开发服务器默认运行在 `5173`，并且已经把 `/api` 代理到 `http://127.0.0.1:8888`。

### 4. 构建前端

```bash
cd web
npm run build
```

这个命令会先执行 TypeScript 检查，再生成生产构建产物。

### 5. 本地访问地址

- 前端：`http://127.0.0.1:5173`
- 后端：`http://127.0.0.1:8888`

## Docker 部署

### 方式一：Docker Compose

```bash
cp .env.example .env
docker compose up -d --build
```

这会启动：

- MySQL
- Redis
- 后端 API
- 前端静态站点

前端通过 Nginx 提供，`/api/` 会转发到后端服务。

### 方式二：部署脚本

Linux 环境下可以直接运行：

```bash
bash deploy-docker.sh
```

脚本提供交互式菜单，支持：

- 全量部署
- 仅重建前端
- 仅重建后端
- 启动、停止、重启
- 查看状态和日志
- Git 更新后重新部署
- 创建默认管理员
- 清理悬空镜像

### 首次部署注意事项

- 首次部署前务必检查 `.env`，尤其是密码、JWT 密钥、SMTP 和 COS 配置
- `databases.sql` 会在 MySQL 首次启动时自动初始化
- 部署脚本支持创建默认管理员，账号为 `admin`，邮箱为 `admin@linux.com`，初始密码为 `123456`
- 登录后请立即修改默认管理员密码

## 接口概览

接口定义来自 `gopan/gopan.api`，运行时路由由 `gopan/internal/handler/routes.go` 注册。

### 公开接口

- `POST /code/send`
- `POST /register`
- `GET /resource/info`
- `GET /user/detail`
- `POST /user/login`

### 登录后接口

- `POST /file/upload`
- `POST /file/preupload`
- `POST /file/chunkupload`
- `POST /file/chunkupload/complete`
- `GET /file/list`
- `DELETE /file/delete`
- `POST /file/rename`
- `PUT /file/move`
- `POST /folder/create`
- `POST /share/create`
- `POST /resource/save`
- `POST /token/refresh`
- `POST /user/repository`

### 管理员接口

以下接口需要管理员权限：

- `GET /admin/overview`
- `GET /admin/users`
- `PUT /admin/user/status`
- `GET /admin/files`
- `DELETE /admin/file`
- `GET /admin/logs`

## 前端接口分层

前端 API 已按功能模块拆分，方便维护：

- `web/src/api/modules/auth.ts`：登录、注册、验证码、用户详情
- `web/src/api/modules/disk.ts`：文件列表、上传、下载、移动、重命名、删除、分享、分片上传
- 其他模块用于后台管理和分享页面逻辑

Axios 统一封装在 `web/src/api/http.ts` 中，包含：

- 自动注入 `Authorization`
- 401 时自动刷新 token
- 刷新失败后清理本地登录态

## 常见工作流

### 文件上传

标准上传和分片上传都支持。前端会先进行预上传判断，利用文件哈希决定是否需要真正传输文件内容。

### 分片上传

大文件会先调用预上传接口，再分片上传到 COS，最后调用分片完成接口提交合并信息。

### 分享文件

分享通过 `share_link` 记录管理，支持设置失效时间，部分逻辑还支持访问密码。

### 后台管理

后台主要做三类事情：

- 查看系统概览
- 管理用户状态和权限
- 查看文件和审计日志

## 开发建议

- 修改接口或配置后，先执行后端编译和前端构建，确保基础链路正常
- 前端登录态保存在浏览器本地存储中，`token` 过期后会自动尝试刷新
- 若启用邮箱验证码或 COS 上传，请先确认对应环境变量已正确配置
- 如果用户已经登录但接口持续返回 `401`，先检查 `user_basic.status`，因为认证中间件会拒绝禁用用户
