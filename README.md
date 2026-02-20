# CHSI Auto Score Query System

一个完整的学信网（CHSI）考研成绩自动查询系统。用户通过前端表单提交个人信息，后端定期自动查询成绩，查询成功后将成绩通过邮件发送给用户。

## 项目概览

```
golang/
├── backend/              # Go后端API服务
│   ├── cmd/             # 应用入口
│   ├── internal/        # 内部逻辑
│   ├── pkg/             # 共享包
│   ├── go.mod          # Go依赖
│   └── README.md       # 后端文档
├── frontend/            # React前端
│   ├── src/            # 源代码
│   ├── public/         # 静态资源
│   └── README.md       # 前端文档
├── docker-compose.yml   # Docker编排
└── README.md           # 本文件
```

## 快速开始

### 使用Docker Compose（推荐）

```bash
# 复制环境变量配置
cp backend/.env.example backend/.env

# 编辑.env文件，填入学信网账户和邮件配置
vim backend/.env

# 启动前后端和数据库
docker-compose up -d
```

访问 `http://localhost:3000` 查看前端界面

### 本地开发

#### 后端

```bash
cd backend
cp .env.example .env
# 编辑.env文件
go mod tidy
go build -o chsi-query ./cmd/server
./chsi-query
```

#### 前端

```bash
cd frontend
npm install
npm run dev
```

## 项目特性

### 后端（Go）

- ✅ 清晰的分层架构（Model-Repository-Service-API）
- ✅ 使用GORM + SQLite数据库
- ✅ 环境变量配置管理
- ✅ 完整的日志系统
- ✅ RESTful API接口
- ✅ 用户隐私保护
  - 用户邮箱存储在数据库，不在配置文件
  - 成功查询后自动删除用户个人信息
- ✅ 防重复提交机制（InfoHash唯一索引）

### 前端（React）

- ✅ 现代化MD3设计
- ✅ 响应式布局
- ✅ 用户表单（姓名、身份证、考生编号、报考单位、邮箱）
- ✅ 表单验证
- ✅ API集成

## API接口

后端服务运行在 `http://localhost:8080`

### 提交个人信息

```
POST /api/submit
Content-Type: application/json

{
  "name": "张三",
  "id_card": "110101199001011234",
  "exam_id": "1100000001",
  "email": "zhangsan@example.com",
  "school_code": "10001"
}
```

### 查询成绩

```
GET /api/score/{email}
```

### 健康检查

```
GET /api/health
```

## 环境变量配置

后端配置文件：`.env`（参考 `backend/.env.example`）

关键变量：
- `CHSI_USERNAME` - 学信网账户
- `CHSI_PASSWORD` - 学信网密码
- `SMTP_SERVER` - 邮件服务器地址
- `SMTP_PORT` - 邮件服务器端口
- `SMTP_USER` - 邮件发送账户
- `SMTP_PASSWORD` - 邮件授权密码
- `DATABASE_DSN` - 数据库路径
- `QUERY_INTERVAL` - 查询间隔（秒）
- `CLEAR_DB_ON_START` - 启动时清空数据库

## 数据库

使用SQLite，数据库文件位置：`./data/chsi.db`

### 用户表结构

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | uint | 主键 |
| name | string | 用户姓名 |
| id_card | string | 身份证号 |
| exam_id | string | 考生编号 |
| email | string | 邮箱 |
| school_code | string | 报考单位代码 |
| info_hash | string | 信息哈希（唯一索引） |
| score | text | 成绩（JSON） |
| notice | text | 通知信息 |
| last_query_at | timestamp | 最后查询时间 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间（软删除） |

## 错误处理

系统处理以下错误情况，并记录详细日志：

1. **登录失败** - 学信网登录不成功
2. **查询失败** - 成绩信息不匹配、成绩未发布等
3. **邮件发送失败** - SMTP连接失败、发送异常等
4. **数据库操作失败** - 连接失败、查询异常等
5. **成功查询** - 成绩查询成功并邮件发送完成

## 日志

日志格式：`[时间戳] [级别]: 消息内容`

支持日志级别：
- DEBUG
- INFO
- WARN
- ERROR

配置日志级别：`LOG_LEVEL=info`

## Docker部署

### dockerfile（后端）

构建后端Docker镜像：

```bash
docker build -f Dockerfile -t chsi-api:latest .
```

### docker-compose.yml

完整的开发和生产部署配置，包括：
- 后端API服务
- 前端服务
- SQLite数据库持久化

## 工作流程

1. 用户在前端填写表单：姓名、身份证、考生编号、报考单位代码、邮箱
2. 前端调用`POST /api/submit`提交数据
3. 后端保存用户信息到数据库（检查重复）
4. 后端定期运行查询任务：
   - 登录学信网（使用配置的学信网账户）
   - 查询用户成绩
   - 解析成绩数据
5. 查询成功后：
   - 发送邮件给用户
   - 删除用户个人信息
   - 记录成功日志
6. 查询失败或邮件发送失败：
   - 发送错误通知（可选）
   - 记录详细错误日志

## 开发指南

### 后端开发

- 在 `internal/service/` 中实现查询逻辑
- 参考 `login.sh` 和 `query.sh` 实现HTTP请求
- 参考 `page.jpg` 实现HTML解析

### 前端开发

- 使用React Hooks
- 集成MD3 Material Design UI库
- 在 `src/components/` 创建可复用组件
- 在 `src/api/` 定义API调用

## 许可证

MIT

## 支持

如有问题，请提交Issue或联系开发者。
