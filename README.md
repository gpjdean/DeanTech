# DeanTech 运维管理平台

DeanTech 是一个面向运维场景的一体化管理平台，覆盖 Kubernetes 集群管理、主机管理、告警中心、日志查询、文档协同、网络连通性测试和系统配置等能力，适合作为企业内部运维控制台。

## 核心能力

- Kubernetes 集群总览、集群列表、节点、Pod、Deployment、Service、ConfigMap、Secret、Ingress、PVC、PV、StorageClass、事件、资源配额和限制范围管理。
- Prometheus 告警、规则、目标监控与静默/抑制规则管理。
- SSL 证书监控中心，支持到期风险、健康度和检查状态展示。
- 主机管理、KVM 管理、SSH 控制台、文件传输、远程命令执行。
- Telnet、Ping、Traceroute、JSON/Base64、MD5、图片压缩、Crontab 配置等工具页。
- SLS 日志查询与配置、OnlyOffice 文档协同、Jenkins 连接与流水线控制。
- 用户管理、个人中心、登录日志、操作日志、系统风格和邮箱配置。

## 技术栈

### 后端

- Go 1.23+
- Gin
- GORM
- MySQL 8.0+
- client-go
- Prometheus HTTP API
- SSH / SFTP / WebSocket

### 前端

- Vue 3
- TypeScript
- Vite
- Element Plus
- Axios
- ECharts
- xterm.js

## 项目结构

```text
.
├── backend/                Go 后端
├── frontend/               Vue 前端
├── deploy/                 部署辅助文件
├── docs/                   详细文档
├── docker-compose.yml      一键部署编排
└── README.md               项目总览
```

## 快速开始

### 本地开发

```bash
# 后端
cd backend
go mod download
go run cmd/main.go

# 前端
cd frontend
npm install
npm run dev
```

默认访问：

- 前端：`http://localhost:80`
- 后端：`http://localhost:8000`

### 默认账号

- 用户名：`admin`
- 密码：`deanit.cn`

## 配置说明

后端主配置位于 [`backend/config/config.yaml`](backend/config/config.yaml)，也支持环境变量覆盖：

| 环境变量 | 说明 |
| --- | --- |
| `SERVER_HOST` | 监听地址 |
| `SERVER_PORT` | 监听端口 |
| `DATABASE_DSN` | MySQL 连接串 |
| `PROMETHEUS_ADDRESS` | Prometheus 地址 |
| `ALERTMANAGER_URL` | Alertmanager 地址 |
| `GIN_MODE` | 运行模式 |

前端生产环境通过 Nginx 转发 `/api` 到后端，开发环境通过 Vite 代理 `/api`。

## 部署方式

- Docker 说明见 [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md)
- SQL 初始化见 [`docs/SQL_INIT.md`](docs/SQL_INIT.md)

## 主要接口

- `/api/dashboard/*`
- `/api/prometheus/*`
- `/api/clusters/*`
- `/api/hosts/*`
- `/api/settings/*`
- `/api/ssl-monitors/*`
- `/api/sls/*`
- `/api/users/*`

## 说明

- 后端启动时会通过 GORM AutoMigrate 自动建表。
- MySQL 是当前项目的主数据库。
- Prometheus 相关页面依赖后端配置中的 `prometheus.addr`，如果目标不可达会显示降级提示。
