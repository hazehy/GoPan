# 论文图表 Mermaid 代码（图2-1 到 图5-3，V2 重制版）

本版设计目标：

1. 学术风格：黑白灰+轻蓝点缀，不花哨。
2. 版面均衡：避免“又长又扁”，优先接近 4:3。
3. 标注清晰：中文业务语义 + 英文技术标识。

## 全局导出建议

1. 优先导出 `SVG`，保证 Word 放大不糊。
2. 常规图导出为 `1600x1200`。
3. 时序图导出为 `1800x1200`。
4. ER 图导出为 `1700x1200`。
5. 图题统一格式：`图X-X 名称`，置于图下。

---

## 图2-1 后端分层与调用关系图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569','secondaryColor':'#EEF2F7'}}}%%
flowchart TB
  subgraph CLIENT[客户端]
    A[浏览器请求]
  end

  subgraph CORE[GoPan 后端核心链路]
    direction TB
    B[路由层 Routes]
    C[处理层 Handler]
    D[业务层 Logic]
    E[数据层 Models]
    B --> C --> D --> E
  end

  subgraph STORE[数据与对象存储]
    direction LR
    F[(MySQL)]
    G[(Redis)]
    H[(腾讯云 COS)]
  end

  A --> B
  E --> F
  D --> G
  D --> H

  M[认证中间件 Auth] -.注入.-> C
  N[管理员中间件 Admin] -.注入.-> C
```

## 图2-2 双令牌续签时序图（1800x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','actorBkg':'#F8FAFC','actorBorder':'#334155','actorTextColor':'#0F172A','signalColor':'#334155','signalTextColor':'#0F172A','lineColor':'#475569','noteBkgColor':'#EEF2F7','noteBorderColor':'#64748B'}}}%%
sequenceDiagram
  participant U as 用户浏览器
  participant B as 业务接口
  participant R as 刷新接口

  U->>B: 携带 AccessToken 请求
  B-->>U: 401（访问令牌失效）
  U->>R: 携带 RefreshToken 刷新

  alt 刷新成功
    R-->>U: 返回新双令牌
    U->>B: 重放原请求
    B-->>U: 200 + 业务数据
  else 刷新失败
    R-->>U: 401/403
    U-->>U: 清理登录态并跳转登录页
  end
```

## 图2-3 对象存储与元数据解耦示意图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  subgraph V[用户目录视图层]
    direction LR
    U1[user_repository\n用户A目录项]
    U2[user_repository\n用户B目录项]
  end

  subgraph P[资源池与对象层]
    direction LR
    R1[repository_pool\nidentity/hash/path]
    O1[(COS Object)]
  end

  U1 -- repository_identity --> R1
  U2 -- repository_identity --> R1
  R1 -- path --> O1

  T[同一物理对象可被多个目录引用\n实现去重与低存储成本]
  R1 -.说明.-> T
```

## 图3-1 系统总体用例图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  U((普通用户))
  A((管理员))

  subgraph UCASE[用户侧核心用例]
    direction LR
    UC1((注册/登录))
    UC2((上传文件))
    UC3((目录管理))
    UC4((重命名/移动/删除))
    UC5((创建分享))
    UC6((保存分享资源))
  end

  subgraph ACASE[管理侧核心用例]
    direction LR
    AC1((数据总览))
    AC2((用户状态管理))
    AC3((全文件治理))
    AC4((审计日志检索))
  end

  U --> UC1
  U --> UC2
  U --> UC3
  U --> UC4
  U --> UC5
  U --> UC6

  A --> AC1
  A --> AC2
  A --> AC3
  A --> AC4
```

## 图3-2 数据流图（DFD）（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  subgraph FE[前端层]
    direction LR
    F1[用户端页面]
    F2[管理员端页面]
  end

  subgraph BE[后端服务层]
    direction LR
    B1[认证与用户服务]
    B2[文件服务]
    B3[分享服务]
    B4[后台与审计服务]
  end

  subgraph DS[数据层]
    direction LR
    DB[(MySQL)]
    RD[(Redis)]
    OS[(腾讯云 COS)]
  end

  F1 -->|登录/注册| B1
  F1 -->|上传/目录| B2
  F1 -->|分享/保存| B3
  F2 -->|统计/治理/日志| B4

  B1 <--> DB
  B1 <--> RD
  B2 <--> DB
  B2 <--> OS
  B3 <--> DB
  B4 <--> DB
```

## 图3-3 普通用户用例图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  U((普通用户))

  subgraph BASE[基础流程]
    direction LR
    A((注册))
    B((登录))
    C((上传文件))
    D((新建文件夹))
  end

  subgraph OPS[文件与分享操作]
    direction LR
    E((文件重命名))
    F((文件移动))
    G((文件删除))
    H((创建分享链接))
    I((保存分享资源))
  end

  C1((秒传判定))
  C2((分片上传))

  U --> A
  U --> B
  U --> C
  U --> D
  U --> E
  U --> F
  U --> G
  U --> H
  U --> I

  C -.包含.-> C1
  C -.包含.-> C2
```

## 图3-4 管理员用例图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  A((管理员))

  subgraph GOV[治理能力]
    direction LR
    U1((查看总览统计))
    U2((查询用户列表))
    U3((启用/禁用用户))
    U4((全文件查询))
    U5((删除违规文件))
  end

  L1((审计日志筛选))
  L2((按日期/动作/用户/扩展名过滤))

  A --> U1
  A --> U2
  A --> U3
  A --> U4
  A --> U5
  A --> L1
  L1 -.包含.-> L2
```

## 图4-1 系统分层架构图（1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  C[浏览器 / Vue3 前端]
  G[REST API 网关]

  subgraph APP[GoPan 应用层]
    direction LR
    H[Handler 处理层]
    L[Logic 业务层]
    M[Model 数据访问层]
  end

  subgraph DATA[数据与存储层]
    direction LR
    D1[(MySQL)]
    D2[(Redis)]
    D3[(腾讯云 COS)]
  end

  C --> G --> H --> L --> M --> D1
  L --> D2
  L --> D3
```

## 图4-2 数据库 E-R 图（1700x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','lineColor':'#475569','primaryTextColor':'#0F172A'}}}%%
erDiagram
  USER_BASIC ||--o{ USER_REPOSITORY : 拥有
  USER_BASIC ||--o{ SHARE_LINK : 创建
  USER_BASIC ||--o{ AUDIT_LOG : 产生

  REPOSITORY_POOL ||--o{ USER_REPOSITORY : 被引用
  REPOSITORY_POOL ||--o{ SHARE_LINK : 被分享

  USER_BASIC {
    string identity PK "用户唯一标识"
    string name "用户名"
    string email "邮箱"
    int role "角色"
    int status "状态"
  }

  REPOSITORY_POOL {
    string identity PK "资源标识"
    string hash UK "文件哈希"
    string name "文件名"
    string ext "扩展名"
    int64 size "大小"
    string path "对象路径"
  }

  USER_REPOSITORY {
    string identity PK "目录项标识"
    int64 parent_id "父目录ID"
    string user_identity FK "所属用户"
    string repository_identity FK "资源标识"
    string name "显示名"
    string ext "扩展名"
  }

  SHARE_LINK {
    string identity PK "分享标识"
    string user_identity FK "分享者"
    string repository_identity FK "分享资源"
    int expires "有效期(天)"
    int click_num "访问次数"
  }

  AUDIT_LOG {
    string identity PK "日志标识"
    string actor_identity FK "操作者标识"
    string actor_name "操作者名称"
    string action "动作类型"
    string target_identity "目标标识"
    string detail "详细描述"
    datetime created_at "创建时间"
  }
```

## 图5-1 登录时序图（1800x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','actorBkg':'#F8FAFC','actorBorder':'#334155','actorTextColor':'#0F172A','signalColor':'#334155','signalTextColor':'#0F172A','lineColor':'#475569'}}}%%
sequenceDiagram
  participant UI as 前端页面
  participant API as 登录接口 /user/login
  participant DB as MySQL
  participant JWT as 令牌服务

  UI->>API: 提交用户名与密码
  API->>DB: 按用户名查询记录
  DB-->>API: 返回用户数据
  API->>API: 校验密码与状态
  API->>JWT: 生成访问令牌与刷新令牌
  JWT-->>API: 返回双令牌
  API-->>UI: 200 + token 信息
  UI-->>UI: 保存登录态
  UI-->>UI: 跳转到 /disk 或 /admin
```

## 图5-2 文件上传流程图（紧凑版，1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  S[选择文件并计算 MD5] --> P[预上传判重 /file/preupload]
  P --> Q{哈希是否存在}

  Q -- 是 --> Y1[直接挂载到用户目录]
  Q -- 否 --> N1[分片上传并完成合并]
  N1 --> N2[写入文件元数据并挂载目录]

  Y1 --> E[上传完成]
  N2 --> E
```

## 图5-3 分享保存流程图（紧凑版，1600x1200）

```mermaid
%%{init: {'theme':'base','themeVariables':{'fontFamily':'Microsoft YaHei, PingFang SC, sans-serif','primaryColor':'#F8FAFC','primaryBorderColor':'#334155','primaryTextColor':'#0F172A','lineColor':'#475569'}}}%%
flowchart TB
  A[创建分享链接并记录 share_link] --> B[访问 /share/:identity 并读取资源]
  B --> C{访问者是否已登录}

  C -- 否 --> L[跳转登录]
  C -- 是 --> D[选择目录并处理同名]
  D --> E[写入 user_repository 并记录审计]
  E --> F[返回保存成功]
```

---

## 最终观感优化建议（Word）

1. 所有图统一宽度：页面宽度的 84%。
2. 图前后段距统一：段前 8pt，段后 10pt。
3. 图题字体：宋体小五，居中。
4. 图与图之间至少空 1 行，避免拥挤。
5. 不要混用彩色截图与黑白流程图，可统一加浅灰边框。
