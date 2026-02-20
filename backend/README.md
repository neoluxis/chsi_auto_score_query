# CHSI Auto Score Query - Backend API

纯后端API服务，使用Go实现。

## 项目结构

```
backend/
├── cmd/
│   └── server/           # 应用程序入口
│       └── main.go       # 主函数
├── internal/
│   ├── api/              # HTTP API层
│   │   ├── server.go     # 服务器初始化和路由注册
│   │   └── handler.go    # HTTP请求处理器
│   ├── db/               # 数据库层
│   │   └── db.go         # GORM初始化
│   ├── logger/           # 日志系统
│   │   └── logger.go     # 日志工具
│   ├── model/            # 数据模型
│   │   └── user.go       # 用户模型
│   ├── repo/             # 数据持久化层
│   │   └── user.go       # 用户仓储
│   └── service/          # 业务逻辑层
│       ├── chsi.go       # CHSI查询服务
│       └── email.go      # 邮件发送服务
├── pkg/
│   ├── config/           # 共享配置
│   │   └── config.go     # 环境变量加载
│   └── utils/            # 工具函数
├── go.mod               # Go模块定义
├── .env.example         # 环境变量示例
├── .gitignore          # Git忽略配置
└── README.md           # 本文件
```

## API端点

- `GET /` - 健康检查
- `GET /api/health` - 服务状态
- `POST /api/submit` - 提交个人信息
  - 请求体：`{"name":"","id_card":"","exam_id":"","email":"","school_code":""}`
- `GET /api/score/{email}` - 查询成绩

## 环境变量配置

复制 `.env.example` 到 `.env`：

```bash
cp .env.example .env
```

配置以下变量：
- `CHSI_USERNAME` - 学信网账户
- `CHSI_PASSWORD` - 学信网密码
- `SMTP_USER` - 邮件发送账户
- `SMTP_PASSWORD` - 邮件授权密码
- 其他配置见 `.env.example`

## 构建与运行

```bash
# 初始化依赖
go mod tidy

# 构建二进制
go build -o chsi-query ./cmd/server

# 运行
./chsi-query
```

## 设计原则

1. **清晰的分层结构**
   - `cmd/` 应用入口
   - `internal/` 内部逻辑
   - `pkg/` 可复用包

2. **数据库隐私保护**
   - 用户邮箱存储在数据库（不在配置文件）
   - 成功查询后自动删除用户个人信息

3. **防重复提交**
   - 使用InfoHash唯一索引预防重复条目

4. **完整的错误处理和日志**
   - 登录失败
   - 查询失败（信息不匹配、成绩未发布等）
   - 邮件发送失败
   - 数据库操作失败

## 后续实现

- [ ] CHSI登录逻辑（参考 `login.sh`）
- [ ] 成绩查询逻辑（参考 `query.sh`）
- [ ] HTML解析逻辑（参考 `page.jpg`）
- [ ] 邮件发送实现（使用SMTP）
- [ ] 后台定时查询调度器
- [ ] 完整的错误处理和日志记录
