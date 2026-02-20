# CHSI Auto Score Query - Frontend

前端Web界面，使用React构建。

## 项目结构

```
frontend/
├── src/
│   ├── components/       # 可复用组件
│   ├── pages/           # 页面组件
│   ├── api/             # API调用接口
│   ├── styles/          # 全局样式
│   └── App.jsx          # 主应用组件
├── public/              # 静态资源
├── package.json         # 项目依赖
└── README.md           # 本文件
```

## 功能

- [x] 用户信息表单（姓名、身份证、考生编号、报考单位代码、邮箱）
- [x] 表单验证
- [ ] 反复查询支持
- [ ] 成绩显示页面
- [ ] 邮件发送后页面提示

## 设计风格

- 现代化 MD3 设计
- 响应式布局
- 优良的用户体验

## API集成

连接到后端服务：
- `POST /api/submit` - 提交个人信息
- `GET /api/score/{email}` - 查询成绩状态

## 环境配置

创建 `.env.local`：

```
VITE_API_BASE_URL=http://localhost:8080
```

## 构建与运行

```bash
# 安装依赖
npm install

# 开发服务器
npm run dev

# 生产构建
npm run build
```
