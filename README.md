# GoPan

GoPan 是一个基于 Go + Vue 3 的网盘系统，支持用户注册登录、文件管理、分片上传、分享链接、管理员后台和审计日志。

## 技术栈

- 后端: Go, go-zero, xorm, MySQL, Redis
- 前端: Vue 3, TypeScript, Vite, Pinia, Vue Router, Axios
- 对象存储: Tencent COS

## 项目结构

```text
GoPan/
  databases.sql                # 数据库初始化脚本
  README.md                    # 项目说明（本文件）
  gopan/                       # 后端服务
    gopan.go                   # 服务入口
    gopan.api                  # API 定义
    etc/gopan-api.yaml         # 默认配置
    etc/gopan-api.local.yaml   # 本地配置（建议）
    internal/                  # handler/logic/svc/types
    helper/                    # 公共辅助能力
    models/                    # 数据模型
  web/                         # 前端工程
    package.json
    vite.config.ts
    src/
```

## 环境要求

- Go: 以 `go.mod` 为准（当前为 `go 1.25.7`）
- Node.js: 建议 `>= 18`
- npm: 建议 `>= 9`
- MySQL: `8.x`
- Redis: `6.x` 及以上

## 快速开始

### Docker 一键部署

如果你在 Linux 服务器上部署，优先使用仓库根目录的 `deploy-docker.sh`。

```bash
chmod +x deploy-docker.sh
./deploy-docker.sh
```

脚本会做这些事：

- 检查并安装 Docker / Docker Compose
- 如果没有 `.env`，自动从 `.env.example` 生成
- 使用 `docker compose up -d --build` 启动整套服务

部署前请先修改 `.env` 里的数据库密码和 `GOPAN_JWT_KEY`，以及 COS / SMTP 相关配置。

### 1. 初始化数据库

1. 创建数据库，例如 `gopan`
2. 执行 `databases.sql`

### 2. 配置后端

后端读取两类配置。

1. YAML 配置文件

- 默认文件: `gopan/etc/gopan-api.yaml`
- 推荐本地文件: `gopan/etc/gopan-api.local.yaml`
- `.gitignore` 已忽略 `gopan/etc/*.local.yaml`

2. 环境变量（敏感配置）

| 变量名                     | 说明          | 默认值             |
| -------------------------- | ------------- | ------------------ |
| `GOPAN_JWT_KEY`            | JWT 密钥      | `change-me-in-env` |
| `GOPAN_FROM_MAIL`          | 发件邮箱      | 空                 |
| `GOPAN_MAIL_PASSWORD`      | 邮箱授权码    | 空                 |
| `GOPAN_SMTP_HOST`          | SMTP 主机     | `smtp.163.com`     |
| `GOPAN_SMTP_PORT`          | SMTP 端口     | `465`              |
| `GOPAN_TENCENT_SECRET_ID`  | COS SecretId  | 空                 |
| `GOPAN_TENCENT_SECRET_KEY` | COS SecretKey | 空                 |
| `GOPAN_COS_BUCKET_URL`     | COS 桶地址    | 空                 |

Windows PowerShell 示例:

```powershell
$env:GOPAN_JWT_KEY="replace-with-strong-random-string"
$env:GOPAN_FROM_MAIL="your_mail@163.com"
$env:GOPAN_MAIL_PASSWORD="your_mail_auth_code"
$env:GOPAN_SMTP_HOST="smtp.163.com"
$env:GOPAN_SMTP_PORT="465"
$env:GOPAN_TENCENT_SECRET_ID="AKIDxxxx"
$env:GOPAN_TENCENT_SECRET_KEY="xxxx"
$env:GOPAN_COS_BUCKET_URL="https://xxx.cos.ap-guangzhou.myqcloud.com"
```

### 3. 启动后端

```bash
cd gopan
go run gopan.go -f etc/gopan-api.local.yaml
```

默认监听: `0.0.0.0:8888`

### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

默认地址: `http://127.0.0.1:5173`

代理规则见 `web/vite.config.ts`: `/api/* -> http://127.0.0.1:8888/*`

## 构建与质量检查

### 后端

```bash
go build ./...
```

### 前端

```bash
cd web
npm run build
```

### 建议测试命令

```bash
go test ./...
```

说明:

- 部分测试依赖外部服务（MySQL/Redis/COS/SMTP）和本地测试资源文件。
- 缺少环境变量时，`mail_test`、`cos_test` 会自动跳过。

## 功能模块

### 用户端

- 登录: `/login`
- 注册: `/register`
- 网盘主页: `/disk`
- 分享页: `/share/:identity`
- 上传增强: 多文件上传队列、任务暂停/继续、断点续传

### 管理端

- 管理页: `/admin`
- 能力: 数据总览、用户管理、文件管理、日志审计

## API 分组

接口定义: `gopan/gopan.api`

运行时路由: `gopan/internal/handler/routes.go`

- 公共接口
  - `POST /user/login`
  - `POST /register`
  - `POST /code/send`
  - `GET /resource/info`
  - `GET /user/detail`
- 用户接口（Auth）
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
- 管理员接口（Auth + Admin）
  - `GET /admin/overview`
  - `GET /admin/users`
  - `PUT /admin/user/status`
  - `GET /admin/files`
  - `DELETE /admin/file`
  - `GET /admin/logs`

## 数据库表

建表脚本: `databases.sql`

- `user_basic`: 用户账户信息
- `repository_pool`: 文件公共池（按哈希去重）
- `user_repository`: 用户文件树
- `share_link`: 分享记录
- `audit_log`: 操作审计日志

## 开发规范

### 后端分层

- `handler`: 参数解析和响应输出
- `logic`: 业务逻辑
- `models`: 数据访问模型
- `types`: 请求/响应结构

### API 变更流程

1. 先修改 `gopan/gopan.api`
2. 再同步 `gopan/internal/types`、`gopan/internal/logic`
3. 最后同步 `web/src/api`、`web/src/types` 与页面调用
4. 更新 README 的接口或行为说明

### 提交前检查

```bash
go fmt ./...
go build ./...
cd web
npm run build
```

### Git 规范（建议）

- 分支: `feature/*`, `fix/*`, `refactor/*`
- 提交信息: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`

## 安全规范

- 不要把真实密钥、邮箱授权码提交到仓库
- 敏感信息统一走环境变量或本地私有配置
- 建议按环境维护配置文件（dev/test/prod）

## Docker 云端部署（Linux）

项目已提供 Docker 化部署文件，可直接在云服务器上通过 `docker compose` 一键启动：

- `docker-compose.yml`
- `gopan/Dockerfile`（后端镜像）
- `gopan/etc/gopan-api.docker.yaml`（后端容器配置模板）
- `web/Dockerfile`（前端镜像）
- `web/nginx/default.conf`（前端 Nginx + API 反向代理）
- `.env.example`（环境变量示例）

### 1. 服务器准备

建议环境：

- Linux x86_64
- Docker >= 24
- Docker Compose Plugin >= 2.20

Ubuntu 示例安装：

```bash
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo $VERSION_CODENAME) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### 2. 拉取代码并配置环境

```bash
git clone <your-repo-url> GoPan
cd GoPan
cp .env.example .env
```

编辑 `.env`，至少修改以下关键配置：

- `MYSQL_ROOT_PASSWORD`
- `GOPAN_JWT_KEY`
- `GOPAN_FROM_MAIL` / `GOPAN_MAIL_PASSWORD`
- `GOPAN_TENCENT_SECRET_ID` / `GOPAN_TENCENT_SECRET_KEY` / `GOPAN_COS_BUCKET_URL`

### 3. 启动服务

```bash
docker compose up -d --build
```

服务说明：

- `mysql`：自动执行 `databases.sql` 初始化表结构
- `redis`：缓存
- `backend`：Go API，容器内监听 `8888`
- `frontend`：Nginx 托管前端并反代 `/api/* -> backend:8888/*`

### 4. 验证部署

```bash
docker compose ps
docker compose logs -f backend
docker compose logs -f frontend
```

浏览器访问：

- `http://<服务器公网IP>/`

API 健康检查示例：

```bash
curl -i http://127.0.0.1:8888/user/detail
```

### 5. 常用运维命令

重启：

```bash
docker compose restart
```

更新发布：

```bash
git pull
docker compose up -d --build
```

停止并保留数据：

```bash
docker compose down
```

彻底清理（会删除 MySQL/Redis 数据卷）：

```bash
docker compose down -v
```

## 常见问题

### 后端启动失败

- 检查 MySQL/Redis 是否可用
- 检查 `gopan/etc/gopan-api.local.yaml` 中连接串
- 先执行 `go build ./...` 验证编译是否通过

### 前端请求 404/跨域

- 确认后端已启动在 `127.0.0.1:8888`
- 检查 `web/vite.config.ts` 代理配置
- 前端请求路径应以 `/api` 开头

### COS 或邮件功能异常

- 检查 `GOPAN_*` 环境变量是否注入
- 检查 COS 桶权限、SMTP 服务和授权码是否有效
