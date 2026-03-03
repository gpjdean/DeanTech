# DeanTech运维管理平台

一个基于Go和Vue3开发的全功能运维管理平台，支持监控告警、K8s集群管理、主机管理、自动化运维等功能，为企业提供高效、可靠的运维解决方案。

## 修复的问题

1. **主机监控资源获取失败**：修复了SSH连接超时问题，添加了3秒超时机制，避免无限加载
2. **Secret/ConfigMap管理页面无数据**：修复了命名空间参数传递和处理逻辑，确保资源正确显示
3. **页面切换不刷新**：修复了路由和组件通信问题，确保菜单切换时页面正确刷新
4. **首页404错误**：修复了Vite代理配置，确保静态资源和API请求正确处理
5. **前端无法连接后端**：修复了API路径配置和代理转发，确保前后端正常通信

## 功能特性

### 1. 监控告警管理
- 实时展示Prometheus告警
- 支持告警状态管理（触发/解决）
- 告警历史查询和统计
- 告警详情查看
- 告警规则管理
- 告警目标监控

### 2. 集群管理
- K8s集群添加和配置
- 集群连接状态测试
- 集群节点管理
- Pod、Deployment、Service管理
- ConfigMap和Secret管理
- 集群资源状态监控

### 3. 主机管理
- 主机添加和配置
- SSH连接测试
- 主机资源监控
- 远程命令执行
- 文件上传下载
- 主机重启和关机
- 主机列表管理

### 4. 实用工具
- 路由跟踪（Traceroute）
- Ping测试
- Telnet测试
- 网络连通性检测

### 5. 告警介质与模板
- 支持多种告警类型（邮件、短信、Webhook等）
- 自定义告警模板
- 模板启用/禁用管理
- 告警介质配置

### 6. 告警抑制和静默
- 静默规则管理
- 抑制规则配置
- 灵活的匹配规则
- 时间范围控制

### 7. 用户与系统管理
- 用户管理（创建、编辑、删除）
- 个人中心设置
- 系统配置管理

## 技术栈

### 后端
- Go 1.21
- Gin Web框架
- GORM ORM框架
- SQLite数据库
- Prometheus客户端库
- Kubernetes客户端库
- SSH客户端库

### 前端
- Vue 3
- TypeScript
- Element Plus UI框架
- Vue Router
- Axios
- Vuex（状态管理）
- ECharts（图表展示）

## 项目结构

```
PromAlert/
├── backend/                 # Go后端代码
│   ├── cmd/                 # 主入口
│   │   └── main.go          # 程序入口
│   ├── api/                 # API路由配置
│   ├── config/              # 配置管理
│   ├── internal/            # 内部包
│   │   └── database/        # 数据库初始化
│   ├── models/              # 数据模型定义
│   ├── services/            # 业务逻辑层
│   │   ├── alert_service.go     # 告警服务
│   │   ├── cluster_service.go   # 集群服务
│   │   ├── host_service.go      # 主机服务
│   │   ├── ssh_service.go       # SSH服务
│   │   └── prometheus_service.go # Prometheus服务
│   └── go.mod               # Go模块定义
├── frontend/                # Vue3前端代码
│   ├── src/                 # 源代码
│   │   ├── assets/          # 静态资源
│   │   ├── components/      # 通用组件
│   │   ├── views/           # 视图页面
│   │   ├── router/          # 路由配置
│   │   ├── api/             # API服务定义
│   │   └── main.ts          # 入口文件
│   ├── index.html           # HTML模板
│   ├── package.json         # 依赖配置
│   └── vite.config.ts       # Vite配置
├── data/                    # 数据目录（数据库文件）
├── docs/                    # 文档
├── public/                  # 公共资源
├── README.md                # 项目说明
└── reset_admin_password.go  # 重置管理员密码工具
```

## 快速开始

### 1. 环境准备

#### 后端依赖
- Go 1.21或更高版本
- SQLite 3（默认数据库）

#### 前端依赖
- Node.js 16或更高版本
- npm 8或更高版本

### 2. 启动后端服务

```bash
cd backend
# 安装依赖
go mod tidy
# 启动服务
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动，所有API都以 `/api` 前缀开头。

### 3. 启动前端服务

```bash
cd frontend
# 安装依赖
npm install
# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

### 4. 访问应用

在浏览器中访问 `http://localhost:3000` 即可使用DeanTech运维管理平台。

默认管理员账号：
- 用户名：admin
- 密码：admin

## 配置说明

### 后端配置

后端支持通过环境变量进行配置，主要配置项：

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| `SERVER_ADDRESS` | `:8080` | 服务器监听地址和端口 |
| `DATABASE_DSN` | `../data/promalert.db` | 数据库连接字符串 |
| `PROMETHEUS_ADDRESS` | `http://localhost:9090` | Prometheus服务器地址 |
| `ALERTMANAGER_URL` | `http://localhost:9093` | AlertManager地址 |

### 前端配置

前端配置主要在 `vite.config.ts` 中，可修改代理配置指向不同的后端地址：

```typescript
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

## API文档

### 集群管理API
- `GET /api/clusters` - 获取集群列表
- `POST /api/clusters` - 创建集群
- `GET /api/clusters/:id` - 获取单个集群
- `PUT /api/clusters/:id` - 更新集群
- `DELETE /api/clusters/:id` - 删除集群
- `POST /api/clusters/:id/test-connection` - 测试集群连接
- `GET /api/clusters/:id/nodes` - 获取集群节点列表
- `GET /api/clusters/:id/pods` - 获取集群Pod列表
- `GET /api/clusters/:id/pods/:name/logs` - 获取Pod日志
- `GET /api/clusters/:id/deployments` - 获取Deployment列表
- `GET /api/clusters/:id/services` - 获取Service列表
- `GET /api/clusters/:id/configmaps` - 获取ConfigMap列表
- `GET /api/clusters/:id/secrets` - 获取Secret列表

### 主机管理API
- `GET /api/hosts` - 获取主机列表
- `POST /api/hosts` - 创建主机
- `GET /api/hosts/:id` - 获取单个主机
- `PUT /api/hosts/:id` - 更新主机
- `DELETE /api/hosts/:id` - 删除主机
- `POST /api/hosts/:id/test-connection` - 测试主机连接
- `GET /api/hosts/:id/stats` - 获取主机资源统计
- `POST /api/hosts/:id/ssh/command` - 执行SSH命令
- `POST /api/hosts/:id/ssh/upload` - 上传文件
- `GET /api/hosts/:id/ssh/download` - 下载文件
- `POST /api/hosts/:id/ssh/restart` - 重启主机
- `POST /api/hosts/:id/ssh/shutdown` - 关机主机

### 告警相关API
- `GET /api/alerts` - 获取告警列表
- `GET /api/alerts/:id` - 获取单个告警
- `PUT /api/alerts/:id/resolve` - 解决告警
- `DELETE /api/alerts/:id` - 删除告警

### 用户管理API
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `GET /api/users/current` - 获取当前用户
- `POST /api/users/register` - 用户注册
- `POST /api/users/login` - 用户登录
- `PUT /api/users/profile` - 更新用户资料
- `POST /api/users/change-password` - 修改密码

### 实用工具API
- `POST /api/tools/traceroute` - 路由跟踪
- `POST /api/tools/ping` - Ping测试
- `POST /api/tools/telnet` - Telnet测试

## 开发说明

### 代码风格

- 后端：遵循Go官方代码风格，使用`go fmt`格式化代码
- 前端：使用TypeScript，遵循ESLint和Prettier规则

### 开发流程

1. 克隆仓库
2. 启动开发环境
3. 编写代码
4. 运行测试
5. 提交代码

### 测试

#### 后端测试

```bash
cd backend
go test ./...
```

#### 前端测试

```bash
cd frontend
npm run test
```

#### 代码检查

```bash
# 后端代码检查
cd backend
go vet ./...

# 前端代码检查
cd frontend
npm run lint
```

## 部署

### 后端部署

1. 编译Go代码

```bash
cd backend
go build -o promalert cmd/main.go
```

2. 运行编译后的二进制文件

```bash
./promalert
```

### 前端部署

1. 构建前端代码

```bash
cd frontend
npm run build
```

2. 将`dist`目录下的文件部署到Web服务器

### Docker部署（推荐）

#### 后端Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o promalert cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/promalert .
COPY --from=builder /app/data ./data
EXPOSE 8080
CMD ["./promalert"]
```

#### 前端Dockerfile

```dockerfile
FROM node:16-alpine AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 安全建议

1. 生产环境中修改默认管理员密码
2. 配置HTTPS加密通信
3. 限制API访问IP
4. 定期备份数据库
5. 启用日志审计

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 项目地址：https://github.com/yourusername/PromAlert
- 邮箱：your.email@example.com
- 开发团队：DeanTech运维团队

## 更新日志

### v1.0.0 (2026-02-08)

**功能特性：**
- ✅ 完整的告警管理系统
- ✅ K8s集群管理
- ✅ 主机管理和SSH操作
- ✅ 实用网络工具
- ✅ 用户管理系统

**修复的问题：**
- ✅ 主机监控资源获取失败
- ✅ Secret/ConfigMap管理页面无数据
- ✅ 页面切换不刷新
- ✅ 首页404错误
- ✅ 前端无法连接后端

**优化：**
- 提高了系统稳定性和性能
- 优化了用户界面和体验
- 增强了错误处理和日志记录
- 完善了API文档和注释

## 致谢

感谢所有为该项目做出贡献的开发者和用户！