# DeanTech 运维管理平台宣传、架构与部署指南

## 1. 文档说明

本文档面向以下场景：

- 对外宣传平台能力与产品价值
- 对内统一介绍平台架构与模块边界
- 为 Docker / Kubernetes 部署提供标准参考
- 为实施、售前、交付和二开团队提供接口索引

相关资料：

- 项目总览见 [README.md](../README.md)
- 基础部署说明见 [DEPLOYMENT.md](./DEPLOYMENT.md)
- 初始化 SQL 说明见 [SQL_INIT.md](./SQL_INIT.md)

---

## 2. 项目简介

DeanTech 是一套面向企业 IT 运维、平台工程、SRE 和基础设施团队的一体化运维管理平台。平台以“统一入口、统一权限、统一操作、统一审计”为核心设计原则，把 Kubernetes 管理、主机管理、KVM 虚拟化、容器运维、告警中心、日志查询、文档协同、批量自动化、数据库运维与常用诊断工具整合到同一套系统中。

它不是单一的 Kubernetes 面板，也不是单一的堡垒机或监控系统，而是更接近企业内部运维中台和基础设施控制台。对外可以作为“企业级一体化运维平台”进行宣传，对内则适合作为平台团队统一交付和标准化运维的核心入口。

### 2.1 产品定位

- 企业级一体化运维管理平台
- 混合云 / 多环境统一控制台
- 基础设施运维中台
- 可审计、可扩展、可集成的运维工作台

### 2.2 核心价值

- 把集群、主机、虚拟化、容器、日志、告警、协同编辑和自动化能力整合进同一平台
- 降低运维团队在多套工具之间切换的成本
- 用菜单权限、角色权限和操作日志建立统一治理体系
- 让团队从“查看状态”走向“执行动作”再到“追踪审计”的闭环管理
- 为平台化建设、标准化交付和自动化运维提供统一底座

### 2.3 适用客户与场景

- 中大型企业内部运维平台建设
- 多套 Kubernetes 集群统一纳管
- 云主机、物理机、虚拟机混合运维
- 面向研发、测试、交付团队提供标准化运维入口
- 运维自动化、批量发布、巡检与诊断
- 内部文档协同、告警联动、审计留痕

---

## 3. 平台卖点

### 3.1 一体化

平台覆盖当前系统菜单中的主要运维场景，包括：

- 平台概览
- 监控告警
- 集群管理
- 主机管理
- 容器管理
- 日志查询
- 实用工具
- 持续集成
- 批量运维
- 协同编辑
- DB 管理
- 娱乐中心
- 系统配置

### 3.2 可操作

DeanTech 不只是展示数据，更强调直接执行运维动作，例如：

- 集群工作负载创建、预检、应用、伸缩、重启、删除
- 主机命令执行、文件上传下载、远程文件编辑、Web SSH
- Docker 容器、镜像、网络、卷、Compose 的在线管理
- KVM 虚拟机的创建、克隆、快照、网卡、磁盘、VNC、RDP 与串口接入
- Redis 连接测试、配置管理、慢日志与 Key 级操作

### 3.3 可审计

- 用户、角色、菜单权限统一管理
- 关键操作自动写入操作日志
- 登录行为单独记录
- 公告、邮箱、系统参数等平台配置可集中维护

### 3.4 可扩展

- 后端以 Go + Gin 服务化组织业务模块
- 前端按业务域拆分页面与菜单
- 外部服务通过配置接入，包括 Prometheus、Alertmanager、OnlyOffice、阿里云 SLS、Jenkins 等
- 可以继续扩展更多云资源、数据库、CMDB、审批、工单或 AI 能力

---

## 4. 功能全景介绍

以下内容基于当前前端菜单定义与后端实际能力整理，适合作为宣传页、产品白皮书或销售资料的功能说明基础。

### 4.1 平台概览

菜单：平台概览

功能说明：

- 汇总平台核心统计指标
- 展示告警、资源、健康状态等全局信息
- 提供统一首页与关键模块入口

宣传表达建议：

“通过统一概览看板，将平台运行态势、资源使用情况与风险状态集中呈现，帮助团队快速掌握全局。”

### 4.2 监控告警

菜单：

- 告警列表
- 告警规则
- 告警介质
- 告警模板
- 静默规则
- SSL 证书监控

功能说明：

- 对接 Prometheus 告警与规则数据
- 支持平台内告警规则管理与手动触发测试
- 支持告警介质配置、测试与模板化通知内容
- 支持静默规则、临时抑制与恢复处理
- 支持 SSL 证书到期、健康度和检测状态监控

适用价值：

- 统一告警入口，减少监控、通知、值班体系分散问题
- 提升告警治理能力，避免噪音告警和通知失控
- 用证书到期监控减少 HTTPS 类基础故障

### 4.3 集群管理

菜单：

- 集群总览
- 集群部署
- 集群列表
- 节点管理
- 资源拓扑
- 应用商店
- 混沌演练
- 命名空间
- Pod
- 工作负载
- 服务
- 配置
- 密钥
- 路由
- 存储卷声明
- 存储卷
- 存储类
- 事件记录
- 资源配额
- 限制范围
- 水平伸缩
- 垂直伸缩
- RBAC

功能说明：

- 支持多集群接入与连接测试
- 支持节点调度、驱逐、标签、污点等节点治理动作
- 支持 Pod 日志、文件、执行、删除等常用运维能力
- 支持 Deployment、StatefulSet、DaemonSet、Job、CronJob 等工作负载管理
- 支持 Service、Ingress、ConfigMap、Secret、PVC、PV、StorageClass 等资源治理
- 支持资源 YAML 查看、预检、应用和删除
- 支持集群级 RBAC 对象管理
- 支持资源拓扑、集群仪表盘、事件查看
- 支持集群部署编排、预检、审批、执行、回滚和诊断

适用价值：

- 统一管理多集群资源，降低日常运维复杂度
- 强化发布、配置、资源治理能力
- 让集群交付、巡检和问题定位更标准化

### 4.4 主机管理

菜单：

- 虚机管理
- KVM 管理

功能说明：

- 支持主机纳管、导入导出、连接测试
- 支持 SSH 命令执行、文件传输、目录浏览、远程文件编辑
- 支持 Web SSH 终端和主机基础状态查看
- 支持 KVM 宿主机管理与虚拟机生命周期管理
- 支持虚拟机创建、克隆、快照、磁盘、网卡和 XML 配置
- 支持 VNC、RDP、串口、SSH 等多种控制台接入

适用价值：

- 让传统主机运维和虚拟化运维纳入统一平台
- 减少在独立 SSH 工具、VNC 客户端、KVM 控制台之间切换

### 4.5 容器管理

菜单：

- 容器列表
- 镜像管理
- 网络管理
- 数据卷管理
- 全局信息
- Compose 管理

功能说明：

- 基于主机维度管理 Docker 环境
- 支持容器创建、启动、停止、重启、删除、日志查看、终端连接
- 支持镜像拉取、删除与构建
- 支持网络、数据卷的创建与删除
- 支持 Compose 项目、服务、日志与 up/down 动作

适用价值：

- 为非 Kubernetes 环境提供统一容器运维能力
- 适用于测试环境、中间件环境和轻量业务环境治理

### 4.6 日志查询

菜单：

- SLS 日志查询
- SLS 配置管理

功能说明：

- 对接阿里云 SLS
- 支持配置管理、连接测试、Project / LogStore 探测
- 支持日志查询、翻页与检索

适用价值：

- 统一纳入云日志平台能力
- 降低运维和研发跨控制台排查问题的成本

### 4.7 实用工具

菜单：

- 路由跟踪
- 端口扫描
- Ping 测试
- 端口测试
- JSON
- Base64
- 图片压缩
- 图床管理
- MD5 加解密
- Crontab 配置
- 子网掩码
- Shell 检测
- Redis 配置
- MySQL 配置

功能说明：

- 提供日常排障与格式处理工具
- 图床管理支持上传配置、资产列表、批量删除、本地回源文件访问
- 图片压缩与 MD5 历史支持持久化
- 网络工具适合快速定位 DNS、端口、链路和连通性问题

适用价值：

- 把运维高频小工具沉淀进平台，避免团队依赖零散网站和本地脚本

### 4.8 持续集成

菜单：Jenkins

功能说明：

- 通过代理接口对接 Jenkins
- 支持统一入口查看和触发持续集成流程

### 4.9 批量运维

菜单：

- 主机列表
- Playbook 管理
- 任务执行
- 可视任务
- 执行历史
- 变量与模板
- 角色库

功能说明：

- 支持主机清单、Playbook、任务、变量模板、角色库管理
- 支持任务执行、失败重试和执行历史追踪
- 支持可视任务工作区与事件流

适用价值：

- 适合批量巡检、批量配置、批量发布和标准化运维执行

### 4.10 协同编辑

菜单：

- OnlyOffice 配置
- 文档管理

功能说明：

- 支持文档上传、下载、删除、批量删除
- 与 OnlyOffice 集成，支持在线协同编辑

适用价值：

- 把 SOP、交付文档、运行手册与平台协作打通

### 4.11 DB 管理

菜单：Redis 管理

功能说明：

- 支持 Redis 实例接入与连接测试
- 支持概览、慢日志、诊断
- 支持配置读取与修改
- 支持 Key 扫描、查看、写入、删除、过期与 DB 统计

### 4.12 娱乐中心

菜单：

- 娱乐中心
- 五子棋
- 2048
- 扫雷

功能说明：

- 作为平台中的轻量互动模块存在
- 不影响平台核心运维定位，可作为团队文化和产品趣味性补充

### 4.13 系统配置

菜单：

- 个人中心
- 用户管理
- 权限管理
- 系统风格
- 邮箱配置
- 公告管理
- 操作日志
- 登录日志

功能说明：

- 用户、角色、菜单权限管理
- 公告发布与撤回
- 系统级参数管理
- 邮件配置与联调测试
- 操作与登录审计

---

## 5. 技术架构

### 5.1 总体架构

DeanTech 采用前后端分离架构，核心组件包括：

- `frontend`：Vue 3 + TypeScript + Element Plus 管理后台
- `backend`：Go + Gin + GORM API 服务
- `mysql`：平台主数据库
- `guacd`：Guacamole 协议代理
- `guac-gateway`：浏览器远程控制网关
- 外部系统：
  - Kubernetes API Server
  - Prometheus
  - Alertmanager
  - 阿里云 SLS
  - OnlyOffice
  - Jenkins
  - 各类图床服务

### 5.2 分层设计

表现层：

- Vue 3 单页应用
- 按业务域划分菜单、路由和页面
- 统一通过 Axios 调用 `/api/*`

接口层：

- Gin Router 统一注册 API
- 每个业务模块由独立 Service 挂载路由
- 中间件负责 CORS、用户状态校验与操作日志记录

业务层：

- `DashboardService`：平台概览数据聚合
- `AlertService` / `MediaService`：告警、介质、模板、静默
- `ClusterService`：集群治理、资源管理、集群部署
- `HostService` / `SSHService`：主机运维与远程执行
- `KVMService`：KVM 与虚拟机管理
- `DockerService`：Docker 容器、镜像、网络、卷、Compose
- `BatchOpsService`：批量运维任务编排与执行
- `DocumentService`：OnlyOffice 文档协同
- `ImageBedService`：图床配置与图片资产管理
- `RedisManagementService`：Redis 运维管理

数据层：

- MySQL 存储平台配置、资产、用户、权限、日志与历史记录
- 启动时通过 GORM AutoMigrate 自动建表

### 5.3 远程控制链路

图形化控制台：

- `frontend` -> `/guac-ws/` -> `guac-gateway` -> `guacd` -> 目标主机 RDP/VNC

文本终端：

- `frontend` -> WebSocket -> `backend` -> SSH / K8s Exec / KVM Serial

### 5.4 权限与审计

- 菜单权限由前端菜单目录和后端权限分组协同控制
- 用户登录后按角色返回权限集
- 页面访问、菜单显示和操作授权统一受控
- 后端操作日志记录关键 API 行为
- 登录日志单独留痕，满足审计和追踪需要

### 5.5 技术栈

后端：

- Go 1.23+
- Gin
- GORM
- MySQL 8.0+
- Kubernetes client-go
- Prometheus HTTP API
- SSH / SFTP
- WebSocket

前端：

- Vue 3
- TypeScript
- Vite
- Element Plus
- Axios
- ECharts
- xterm.js
- noVNC / Guacamole Web 相关组件

---

## 6. 容器环境部署文档

### 6.1 组件组成

推荐使用 Docker Compose 部署以下服务：

- `mysql`
- `backend`
- `frontend`
- `guacd`
- `guac-gateway`

其中：

- `frontend` 提供 Web 访问入口
- `backend` 提供业务 API
- `mysql` 存储业务数据
- `guacd` 与 `guac-gateway` 提供浏览器远程控制能力
- Prometheus、Alertmanager、OnlyOffice、SLS、Jenkins 可复用企业已有服务

### 6.2 环境要求

- Docker 24+
- Docker Compose v2
- 至少 2 核 CPU、4 GB 内存用于基础体验
- 可访问 MySQL、Prometheus、Alertmanager 或其他外部依赖

### 6.3 镜像构建

```bash
docker build -t deantech-backend ./backend
docker build -t deantech-frontend ./frontend
docker build -t deantech-guac-gateway ./guac-gateway
```

### 6.4 Compose 启动

项目已提供 [docker-compose.yml](../docker-compose.yml)。

启动：

```bash
docker compose up -d --build
```

停止：

```bash
docker compose down
```

查看日志：

```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
docker compose logs -f guac-gateway
```

### 6.5 默认端口

| 服务 | 默认端口 | 说明 |
| --- | --- | --- |
| `frontend` | `80` | 平台 Web 入口 |
| `backend` | `8000` | 后端 API 服务 |
| `mysql` | `3306` | 主数据库 |
| `guacd` | `4822` | Guacamole 代理 |
| `guac-gateway` | `8081` | 远程桌面网关 |

### 6.6 关键环境变量

| 变量 | 说明 |
| --- | --- |
| `GIN_MODE` | Go 服务运行模式 |
| `SERVER_HOST` | 监听地址 |
| `SERVER_PORT` | 监听端口 |
| `DATABASE_DSN` | MySQL 连接串 |
| `PROMETHEUS_ADDRESS` | Prometheus 地址 |
| `ALERTMANAGER_URL` | Alertmanager 地址 |
| `GUACAMOLE_WS_URL` | 前端访问远程控制台路径 |
| `GUACAMOLE_TOKEN_KEY` | Guacamole 网关令牌密钥 |

### 6.7 生产部署建议

- MySQL 建议使用独立数据库服务，不建议长期使用默认本地容器数据
- 前端建议由企业 Nginx / Ingress 统一接入 HTTPS
- `backend` 建议通过环境变量或配置中心注入参数
- `backend` 与 `guac-gateway` 必须使用相同的 `GUACAMOLE_TOKEN_KEY`
- 如果 Prometheus、Alertmanager 不在同一网络，应改成实际可达地址
- 建议单独监控 `backend`、`guac-gateway` 与 `guacd` 的健康状态

---

## 7. Kubernetes 环境部署文档

### 7.1 推荐资源拓扑

- Namespace：`deantech`
- Secret：数据库、令牌、外部服务密钥
- ConfigMap：后端配置、前端 Nginx 配置
- Deployment：
  - `deantech-backend`
  - `deantech-frontend`
  - `deantech-guacd`
  - `deantech-guac-gateway`
- Service：
  - `deantech-backend`
  - `deantech-frontend`
  - `deantech-guacd`
  - `deantech-guac-gateway`
- Ingress：统一暴露前端域名、`/api` 和 `/guac-ws/`

### 7.2 部署顺序

1. 创建 Namespace
2. 创建 Secret / ConfigMap
3. 准备 MySQL 或外部数据库
4. 部署 `backend`
5. 部署 `guacd` 与 `guac-gateway`
6. 部署 `frontend`
7. 配置 Service 与 Ingress

### 7.3 示例 Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: deantech-secret
  namespace: deantech
type: Opaque
stringData:
  DATABASE_DSN: deantech:deantech123@tcp(mysql:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local
  GUACAMOLE_TOKEN_KEY: deantech-guacamole-secret
```

### 7.4 示例 Backend Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-backend
  namespace: deantech
spec:
  replicas: 2
  selector:
    matchLabels:
      app: deantech-backend
  template:
    metadata:
      labels:
        app: deantech-backend
    spec:
      containers:
        - name: backend
          image: deantech-backend:latest
          ports:
            - containerPort: 8000
          env:
            - name: GIN_MODE
              value: release
            - name: SERVER_HOST
              value: 0.0.0.0
            - name: SERVER_PORT
              value: "8000"
            - name: DATABASE_DSN
              valueFrom:
                secretKeyRef:
                  name: deantech-secret
                  key: DATABASE_DSN
            - name: PROMETHEUS_ADDRESS
              value: http://prometheus.monitoring.svc.cluster.local:9090
            - name: ALERTMANAGER_URL
              value: http://alertmanager.monitoring.svc.cluster.local:9093
            - name: GUACAMOLE_WS_URL
              value: /guac-ws/
            - name: GUACAMOLE_TOKEN_KEY
              valueFrom:
                secretKeyRef:
                  name: deantech-secret
                  key: GUACAMOLE_TOKEN_KEY
          readinessProbe:
            httpGet:
              path: /health
              port: 8000
          livenessProbe:
            httpGet:
              path: /health
              port: 8000
```

### 7.5 示例 Frontend Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-frontend
  namespace: deantech
spec:
  replicas: 2
  selector:
    matchLabels:
      app: deantech-frontend
  template:
    metadata:
      labels:
        app: deantech-frontend
    spec:
      containers:
        - name: frontend
          image: deantech-frontend:latest
          ports:
            - containerPort: 80
```

### 7.6 示例 Guacamole 组件

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-guacd
  namespace: deantech
spec:
  replicas: 1
  selector:
    matchLabels:
      app: deantech-guacd
  template:
    metadata:
      labels:
        app: deantech-guacd
    spec:
      containers:
        - name: guacd
          image: docker.m.daocloud.io/guacamole/guacd:1.5.5
          ports:
            - containerPort: 4822
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-guac-gateway
  namespace: deantech
spec:
  replicas: 1
  selector:
    matchLabels:
      app: deantech-guac-gateway
  template:
    metadata:
      labels:
        app: deantech-guac-gateway
    spec:
      containers:
        - name: guac-gateway
          image: deantech-guac-gateway:latest
          ports:
            - containerPort: 8081
          env:
            - name: PORT
              value: "8081"
            - name: GUACD_HOST
              value: deantech-guacd
            - name: GUACD_PORT
              value: "4822"
            - name: GUACAMOLE_TOKEN_KEY
              valueFrom:
                secretKeyRef:
                  name: deantech-secret
                  key: GUACAMOLE_TOKEN_KEY
```

### 7.7 示例 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: deantech
  namespace: deantech
spec:
  ingressClassName: nginx
  rules:
    - host: deantech.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: deantech-frontend
                port:
                  number: 80
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: deantech-backend
                port:
                  number: 8000
          - path: /guac-ws/
            pathType: Prefix
            backend:
              service:
                name: deantech-guac-gateway
                port:
                  number: 8081
```

### 7.8 K8s 生产建议

- 为前后端分别配置 HPA 与 PDB
- 敏感配置统一使用 Secret 托管
- Ingress 需要开启 WebSocket 与较长超时
- 数据库建议使用外部高可用数据库
- 对 `backend`、`guac-gateway` 增加错误率、连接数与响应时间监控
- 使用独立域名和 HTTPS 证书对外暴露平台入口

---

## 8. 功能接口设计说明

### 8.1 通用约定

- API 前缀：`/api`
- 健康检查：`/health`
- 返回格式：JSON 为主
- 文件上传：`multipart/form-data`
- WebSocket：用于 SSH、K8s Exec、Docker 终端、KVM 串口/VNC 等实时会话

### 8.2 鉴权与权限

当前平台的用户会话依赖用户信息与权限分组控制，前端会根据菜单权限渲染页面，后端在路由侧统一接入用户状态校验和操作日志中间件。

核心能力：

- 用户登录、注册、找回密码
- 获取当前用户
- 用户资料更新、密码修改
- 用户管理
- 权限组与权限目录管理

### 8.3 接口使用建议

对外宣传时不建议把所有接口都当作开放 API 产品能力，而应按“平台内部能力接口”进行介绍。对实施、集成或二开团队，可以将以下接口分组作为交付手册的重点。

---

## 9. 重点业务接口总览

以下接口均已根据当前后端真实路由校对，适合作为“有用接口”说明文档。

### 9.1 平台概览与监控

Dashboard：

- `GET /api/dashboard/stats`
- `GET /api/dashboard/resources`
- `GET /api/dashboard/system-info`
- `GET /api/dashboard/health-status`

Prometheus：

- `GET /api/prometheus/alerts`
- `GET /api/prometheus/rules`
- `GET /api/prometheus/targets`

告警中心：

- `GET /api/alerts`
- `GET /api/alerts/:id`
- `PUT /api/alerts/:id/resolve`
- `DELETE /api/alerts/:id`
- `DELETE /api/alerts/:id/suppression`

告警规则：

- `GET /api/alert-rules`
- `POST /api/alert-rules`
- `GET /api/alert-rules/:id`
- `PUT /api/alert-rules/:id`
- `POST /api/alert-rules/:id/run`
- `DELETE /api/alert-rules/:id`

静默规则：

- `GET /api/silences`
- `POST /api/silences`
- `POST /api/silences/batch-delete`
- `GET /api/silences/:id`
- `DELETE /api/silences/:id`

告警介质与模板：

- `GET /api/media`
- `POST /api/media`
- `POST /api/media/test`
- `GET /api/media/:id`
- `POST /api/media/:id/test`
- `PUT /api/media/:id`
- `DELETE /api/media/:id`
- `GET /api/templates`
- `POST /api/templates`
- `GET /api/templates/:id`
- `PUT /api/templates/:id`
- `DELETE /api/templates/:id`

SSL 证书监控：

- `GET /api/ssl-monitors`
- `POST /api/ssl-monitors`
- `PUT /api/ssl-monitors/:id`
- `DELETE /api/ssl-monitors/:id`
- `POST /api/ssl-monitors/:id/check`
- `POST /api/ssl-monitors/:id/test-notify`
- `POST /api/ssl-monitors/check-all`

### 9.2 用户、权限与系统配置

用户与登录：

- `POST /api/users/login`
- `POST /api/users/register`
- `GET /api/users/current`
- `PUT /api/users/profile`
- `POST /api/users/change-password`
- `POST /api/users/forgot-password/verify-email`
- `POST /api/users/forgot-password/reset`

用户管理：

- `GET /api/users`
- `POST /api/users`
- `GET /api/users/:id`
- `PUT /api/users/:id`
- `DELETE /api/users/:id`

权限管理：

- `GET /api/permissions`
- `GET /api/permissions/catalog`
- `POST /api/permissions`
- `PUT /api/permissions/:id`
- `DELETE /api/permissions/:id`

系统参数：

- `GET /api/settings`
- `GET /api/settings/:key`
- `POST /api/settings`
- `PUT /api/settings/:key`
- `DELETE /api/settings/:key`

邮箱：

- `POST /api/email/test`
- `GET /api/email-settings`
- `PUT /api/email-settings`

公告：

- `GET /api/announcements`
- `GET /api/announcements/published`
- `GET /api/announcements/:id`
- `POST /api/announcements`
- `PUT /api/announcements/:id`
- `DELETE /api/announcements/:id`
- `POST /api/announcements/:id/publish`
- `POST /api/announcements/:id/unpublish`

审计日志：

- `GET /api/logs/operation`
- `GET /api/logs/operation-filters`
- `GET /api/logs/operation/:id`
- `DELETE /api/logs/operation/:id`
- `DELETE /api/logs/operation`
- `GET /api/logs/login`
- `GET /api/logs/login/:id`
- `DELETE /api/logs/login/:id`
- `DELETE /api/logs/login`

### 9.3 集群管理

集群与仪表盘：

- `GET /api/clusters`
- `POST /api/clusters`
- `GET /api/clusters/:id`
- `PUT /api/clusters/:id`
- `DELETE /api/clusters/:id`
- `POST /api/clusters/:id/test-connection`
- `GET /api/clusters/:id/dashboard`
- `GET /api/clusters/:id/resource-topology`

节点管理：

- `GET /api/clusters/:id/nodes`
- `PUT /api/clusters/:id/nodes/:name/schedule`
- `POST /api/clusters/:id/nodes/:name/drain`
- `PUT /api/clusters/:id/nodes/:name/taints`
- `PUT /api/clusters/:id/nodes/:name/labels`
- `DELETE /api/clusters/:id/nodes/:name`

Pod 运维：

- `GET /api/clusters/:id/pods`
- `GET /api/clusters/:id/pods/:name/logs`
- `DELETE /api/clusters/:id/pods/:name`
- `DELETE /api/clusters/:id/pods`
- `GET /api/clusters/:id/pods/:name/files`
- `POST /api/clusters/:id/pods/:name/files/upload`
- `GET /api/clusters/:id/pods/:name/files/download`
- `GET /api/clusters/:id/pods/:name/exec`

工作负载：

- `GET /api/clusters/:id/deployments`
- `POST /api/clusters/:id/deployments`
- `POST /api/clusters/:id/deployments/precheck`
- `PUT /api/clusters/:id/deployments/:name/scale`
- `POST /api/clusters/:id/deployments/:name/restart`
- `DELETE /api/clusters/:id/deployments/:name`
- `GET /api/clusters/:id/statefulsets`
- `DELETE /api/clusters/:id/statefulsets/:name`
- `GET /api/clusters/:id/daemonsets`
- `DELETE /api/clusters/:id/daemonsets/:name`
- `GET /api/clusters/:id/jobs`
- `DELETE /api/clusters/:id/jobs/:name`
- `GET /api/clusters/:id/cronjobs`
- `DELETE /api/clusters/:id/cronjobs/:name`
- `GET /api/clusters/:id/workloads/:type/:name/yaml`
- `DELETE /api/clusters/:id/workloads/:type/:name`
- `POST /api/clusters/:id/workloads/precheck`
- `PUT /api/clusters/:id/workloads/apply`
- `POST /api/clusters/:id/workloads/apply`

服务与配置资源：

- `GET /api/clusters/:id/services`
- `DELETE /api/clusters/:id/services/:name`
- `GET /api/clusters/:id/services/:name/yaml`
- `POST /api/clusters/:id/services/precheck`
- `PUT /api/clusters/:id/services/apply`
- `POST /api/clusters/:id/services/apply`
- `GET /api/clusters/:id/configmaps`
- `GET /api/clusters/:id/configmaps/:name/yaml`
- `POST /api/clusters/:id/configmaps/precheck`
- `PUT /api/clusters/:id/configmaps/apply`
- `POST /api/clusters/:id/configmaps/apply`
- `GET /api/clusters/:id/secrets`
- `DELETE /api/clusters/:id/secrets/:name`
- `GET /api/clusters/:id/secrets/:name/yaml`
- `POST /api/clusters/:id/secrets/precheck`
- `PUT /api/clusters/:id/secrets/apply`
- `POST /api/clusters/:id/secrets/apply`
- `GET /api/clusters/:id/ingresses`
- `DELETE /api/clusters/:id/ingresses/:name`
- `GET /api/clusters/:id/ingresses/:name/yaml`
- `POST /api/clusters/:id/ingresses/precheck`
- `PUT /api/clusters/:id/ingresses/apply`
- `POST /api/clusters/:id/ingresses/apply`

命名空间、存储与配额：

- `GET /api/clusters/:id/namespaces`
- `POST /api/clusters/:id/namespaces`
- `DELETE /api/clusters/:id/namespaces/:name`
- `GET /api/clusters/:id/pvcs`
- `DELETE /api/clusters/:id/pvcs/:name`
- `GET /api/clusters/:id/pvcs/:name/yaml`
- `POST /api/clusters/:id/pvcs/precheck`
- `PUT /api/clusters/:id/pvcs/apply`
- `POST /api/clusters/:id/pvcs/apply`
- `GET /api/clusters/:id/pvs`
- `DELETE /api/clusters/:id/pvs/:name`
- `GET /api/clusters/:id/pvs/:name/yaml`
- `POST /api/clusters/:id/pvs/precheck`
- `PUT /api/clusters/:id/pvs/apply`
- `POST /api/clusters/:id/pvs/apply`
- `GET /api/clusters/:id/storageclasses`
- `GET /api/clusters/:id/storageclasses/:name/yaml`
- `DELETE /api/clusters/:id/storageclasses/:name`
- `POST /api/clusters/:id/storageclasses/precheck`
- `PUT /api/clusters/:id/storageclasses/apply`
- `POST /api/clusters/:id/storageclasses/apply`
- `GET /api/clusters/:id/resourcequotas`
- `GET /api/clusters/:id/resourcequotas/:name/yaml`
- `DELETE /api/clusters/:id/resourcequotas/:name`
- `POST /api/clusters/:id/resourcequotas/precheck`
- `PUT /api/clusters/:id/resourcequotas/apply`
- `POST /api/clusters/:id/resourcequotas/apply`
- `GET /api/clusters/:id/limitranges`
- `GET /api/clusters/:id/limitranges/:name/yaml`
- `DELETE /api/clusters/:id/limitranges/:name`
- `POST /api/clusters/:id/limitranges/precheck`
- `PUT /api/clusters/:id/limitranges/apply`
- `POST /api/clusters/:id/limitranges/apply`

自动伸缩、事件与混沌演练：

- `GET /api/clusters/:id/horizontalautoscalers`
- `GET /api/clusters/:id/horizontalautoscalers/:name/yaml`
- `DELETE /api/clusters/:id/horizontalautoscalers/:name`
- `POST /api/clusters/:id/horizontalautoscalers/precheck`
- `PUT /api/clusters/:id/horizontalautoscalers/apply`
- `POST /api/clusters/:id/horizontalautoscalers/apply`
- `GET /api/clusters/:id/verticalautoscalers`
- `GET /api/clusters/:id/verticalautoscalers/capability`
- `GET /api/clusters/:id/verticalautoscalers/:name/yaml`
- `DELETE /api/clusters/:id/verticalautoscalers/:name`
- `POST /api/clusters/:id/verticalautoscalers/precheck`
- `PUT /api/clusters/:id/verticalautoscalers/apply`
- `POST /api/clusters/:id/verticalautoscalers/apply`
- `GET /api/clusters/:id/events`
- `GET /api/clusters/:id/chaos/experiments`
- `POST /api/clusters/:id/chaos/experiments`
- `DELETE /api/clusters/:id/chaos/experiments`

RBAC：

- `GET /api/clusters/:id/rbac/cluster-roles`
- `GET /api/clusters/:id/rbac/cluster-roles/:name/yaml`
- `POST /api/clusters/:id/rbac/cluster-roles`
- `PUT /api/clusters/:id/rbac/cluster-roles/:name`
- `DELETE /api/clusters/:id/rbac/cluster-roles/:name`
- `GET /api/clusters/:id/rbac/roles`
- `GET /api/clusters/:id/rbac/roles/:name/yaml`
- `POST /api/clusters/:id/rbac/roles`
- `PUT /api/clusters/:id/rbac/roles/:name`
- `DELETE /api/clusters/:id/rbac/roles/:name`
- `GET /api/clusters/:id/rbac/cluster-role-bindings`
- `GET /api/clusters/:id/rbac/cluster-role-bindings/:name/yaml`
- `POST /api/clusters/:id/rbac/cluster-role-bindings`
- `PUT /api/clusters/:id/rbac/cluster-role-bindings/:name`
- `DELETE /api/clusters/:id/rbac/cluster-role-bindings/:name`
- `GET /api/clusters/:id/rbac/role-bindings`
- `GET /api/clusters/:id/rbac/role-bindings/:name/yaml`
- `POST /api/clusters/:id/rbac/role-bindings`
- `PUT /api/clusters/:id/rbac/role-bindings/:name`
- `DELETE /api/clusters/:id/rbac/role-bindings/:name`
- `GET /api/clusters/:id/rbac/service-accounts`
- `GET /api/clusters/:id/rbac/service-accounts/:name/yaml`
- `POST /api/clusters/:id/rbac/service-accounts`
- `PUT /api/clusters/:id/rbac/service-accounts/:name`
- `DELETE /api/clusters/:id/rbac/service-accounts/:name`

集群部署编排：

- `GET /api/cluster-deployments/options`
- `GET /api/cluster-deployments`
- `GET /api/cluster-deployments/:id`
- `GET /api/cluster-deployments/:id/revisions`
- `POST /api/cluster-deployments`
- `PUT /api/cluster-deployments/:id`
- `DELETE /api/cluster-deployments/:id`
- `POST /api/cluster-deployments/preview`
- `POST /api/cluster-deployments/precheck`
- `POST /api/cluster-deployments/:id/precheck`
- `POST /api/cluster-deployments/:id/diagnose`
- `GET /api/cluster-deployments/:id/diagnostics/export`
- `POST /api/cluster-deployments/:id/approve`
- `POST /api/cluster-deployments/:id/reject`
- `POST /api/cluster-deployments/:id/run`
- `POST /api/cluster-deployments/:id/retry`
- `POST /api/cluster-deployments/:id/rollback`
- `GET /api/cluster-deployments/:id/events`

### 9.4 主机、容器与虚拟化

主机管理：

- `GET /api/hosts`
- `GET /api/hosts/export`
- `POST /api/hosts/import`
- `POST /api/hosts/batch-delete`
- `POST /api/hosts`
- `GET /api/hosts/:id`
- `GET /api/hosts/:id/console-config`
- `PUT /api/hosts/:id`
- `DELETE /api/hosts/:id`
- `POST /api/hosts/:id/test-connection`

SSH 与文件管理：

- `POST /api/hosts/:id/ssh/command`
- `POST /api/hosts/:id/ssh/upload`
- `GET /api/hosts/:id/ssh/download`
- `POST /api/hosts/:id/ssh/restart`
- `POST /api/hosts/:id/ssh/shutdown`
- `GET /api/hosts/:id/ssh/files`
- `GET /api/hosts/:id/ssh/file-content`
- `PUT /api/hosts/:id/ssh/file-content`
- `POST /api/hosts/:id/ssh/mkdir`
- `GET /api/hosts/:id/ssh/ws`
- `GET /api/hosts/:id/stats`

Docker 运维：

- `GET /api/hosts/:id/docker/overview`
- `GET /api/hosts/:id/docker/containers`
- `POST /api/hosts/:id/docker/containers`
- `GET /api/hosts/:id/docker/containers/:container/terminal/ws`
- `POST /api/hosts/:id/docker/containers/:container/start`
- `POST /api/hosts/:id/docker/containers/:container/stop`
- `POST /api/hosts/:id/docker/containers/:container/restart`
- `DELETE /api/hosts/:id/docker/containers/:container`
- `GET /api/hosts/:id/docker/containers/:container/logs`
- `GET /api/hosts/:id/docker/containers/:container/inspect`
- `GET /api/hosts/:id/docker/images`
- `POST /api/hosts/:id/docker/images/pull`
- `DELETE /api/hosts/:id/docker/images`
- `POST /api/hosts/:id/docker/images/build`
- `GET /api/hosts/:id/docker/networks`
- `POST /api/hosts/:id/docker/networks`
- `DELETE /api/hosts/:id/docker/networks/:network`
- `GET /api/hosts/:id/docker/volumes`
- `POST /api/hosts/:id/docker/volumes`
- `DELETE /api/hosts/:id/docker/volumes/:volume`
- `GET /api/hosts/:id/docker/compose/projects`
- `GET /api/hosts/:id/docker/compose/services`
- `GET /api/hosts/:id/docker/compose/logs`
- `POST /api/hosts/:id/docker/compose/service-action`
- `POST /api/hosts/:id/docker/compose/up`
- `POST /api/hosts/:id/docker/compose/down`

KVM 与虚拟机：

- `GET /api/kvm/hosts`
- `POST /api/kvm/hosts`
- `GET /api/kvm/hosts/:id`
- `PUT /api/kvm/hosts/:id`
- `DELETE /api/kvm/hosts/:id`
- `POST /api/kvm/hosts/:id/test`
- `POST /api/kvm/hosts/:id/vm-control`
- `GET /api/kvm/hosts/:id/vms`
- `POST /api/kvm/hosts/:id/vms`
- `GET /api/kvm/hosts/:id/vms/:vmId/detail`
- `GET /api/kvm/hosts/:id/image-files`
- `GET /api/kvm/hosts/:id/image-catalog`
- `GET /api/kvm/hosts/:id/network-sources`
- `GET /api/kvm/hosts/:id/create-options`
- `POST /api/kvm/hosts/:id/image-download`
- `PUT /api/kvm/hosts/:id/vms/:vmId/manual-ip`
- `GET /api/kvm/hosts/:id/vms/:vmId/xml`
- `PUT /api/kvm/hosts/:id/vms/:vmId/xml`
- `PUT /api/kvm/hosts/:id/vms/:vmId/config`
- `POST /api/kvm/hosts/:id/vms/:vmId/disks`
- `DELETE /api/kvm/hosts/:id/vms/:vmId/disks/:target`
- `POST /api/kvm/hosts/:id/vms/:vmId/interfaces`
- `DELETE /api/kvm/hosts/:id/vms/:vmId/interfaces/:mac`
- `POST /api/kvm/hosts/:id/vms/:vmId/control`
- `POST /api/kvm/hosts/:id/vms/:vmId/clone`
- `POST /api/kvm/hosts/:id/vms/:vmId/snapshots`
- `POST /api/kvm/hosts/:id/vms/:vmId/snapshots/:snapshot/revert`
- `DELETE /api/kvm/hosts/:id/vms/:vmId/snapshots/:snapshot`
- `GET /api/kvm/hosts/:id/vms/:vmId/serial/ws`
- `GET /api/kvm/hosts/:id/vms/:vmId/ssh/ws`
- `GET /api/kvm/hosts/:id/vms/:vmId/vnc-info`
- `GET /api/kvm/hosts/:id/vms/:vmId/vnc-password`
- `POST /api/kvm/hosts/:id/vms/:vmId/vnc-password`
- `GET /api/kvm/hosts/:id/vms/:vmId/vnc/ws`
- `POST /api/kvm/hosts/:id/vms/:vmId/rdp/session`
- `POST /api/kvm/hosts/:id/vms/:vmId/start`
- `POST /api/kvm/hosts/:id/vms/:vmId/shutdown`
- `POST /api/kvm/hosts/:id/vms/:vmId/reboot`
- `POST /api/kvm/hosts/:id/vms/:vmId/reset`
- `POST /api/kvm/hosts/:id/vms/:vmId/suspend`
- `POST /api/kvm/hosts/:id/vms/:vmId/resume`
- `POST /api/kvm/hosts/:id/vms/:vmId/destroy`

### 9.5 文档、图床与批量运维

OnlyOffice 配置：

- `GET /api/onlyoffice/config`
- `POST /api/onlyoffice/config`
- `POST /api/onlyoffice/test-connection`

文档管理：

- `POST /api/documents/upload`
- `GET /api/documents`
- `GET /api/documents/:id`
- `DELETE /api/documents/:id`
- `POST /api/documents/batch-delete`
- `GET /api/documents/:id/download`
- `POST /api/documents/:id/callback`

图床管理：

- `GET /api/image-bed/provider-config`
- `PUT /api/image-bed/provider-config`
- `POST /api/image-bed/provider-config/test`
- `GET /api/image-bed/assets`
- `POST /api/image-bed/assets`
- `POST /api/image-bed/assets/upload`
- `PUT /api/image-bed/assets/:id`
- `DELETE /api/image-bed/assets/:id`
- `POST /api/image-bed/assets/batch-delete`
- `GET /api/image-bed/assets/:id/file`

批量运维：

- `GET /api/batch-ops/playbooks`
- `POST /api/batch-ops/playbooks`
- `GET /api/batch-ops/playbooks/:id`
- `PUT /api/batch-ops/playbooks/:id`
- `DELETE /api/batch-ops/playbooks/:id`
- `GET /api/batch-ops/host-lists`
- `POST /api/batch-ops/host-lists`
- `GET /api/batch-ops/host-lists/:id`
- `PUT /api/batch-ops/host-lists/:id`
- `DELETE /api/batch-ops/host-lists/:id`
- `GET /api/batch-ops/tasks`
- `POST /api/batch-ops/tasks`
- `GET /api/batch-ops/tasks/:id`
- `PUT /api/batch-ops/tasks/:id`
- `DELETE /api/batch-ops/tasks/:id`
- `POST /api/batch-ops/tasks/:id/execute`
- `GET /api/batch-ops/executions`
- `GET /api/batch-ops/executions/:id`
- `DELETE /api/batch-ops/executions/:id`
- `POST /api/batch-ops/executions/batch-delete`
- `POST /api/batch-ops/executions/:id/retry`
- `GET /api/batch-ops/variable-templates`
- `POST /api/batch-ops/variable-templates`
- `GET /api/batch-ops/variable-templates/:id`
- `PUT /api/batch-ops/variable-templates/:id`
- `DELETE /api/batch-ops/variable-templates/:id`
- `GET /api/batch-ops/roles`
- `POST /api/batch-ops/roles`
- `GET /api/batch-ops/roles/:id`
- `PUT /api/batch-ops/roles/:id`
- `DELETE /api/batch-ops/roles/:id`
- `GET /api/batch-ops/visual-tasks/workspace`
- `POST /api/batch-ops/visual-tasks/launch`
- `GET /api/batch-ops/visual-tasks/executions/:id/events`

### 9.6 日志、数据库与工具

SLS：

- `GET /api/sls/configs`
- `GET /api/sls/configs/:id`
- `POST /api/sls/configs`
- `PUT /api/sls/configs/:id`
- `DELETE /api/sls/configs/:id`
- `POST /api/sls/configs/:id/test-connection`
- `GET /api/sls/configs/:id/projects`
- `GET /api/sls/configs/:id/logstores`
- `POST /api/sls/query`

Redis 管理：

- `GET /api/redis-management/instances`
- `GET /api/redis-management/instances/:id`
- `POST /api/redis-management/instances`
- `PUT /api/redis-management/instances/:id`
- `DELETE /api/redis-management/instances/:id`
- `POST /api/redis-management/test`
- `POST /api/redis-management/overview`
- `POST /api/redis-management/slowlog`
- `POST /api/redis-management/diagnostics`
- `POST /api/redis-management/config/get`
- `POST /api/redis-management/config/set`
- `POST /api/redis-management/keys/scan`
- `POST /api/redis-management/keys/get`
- `POST /api/redis-management/keys/set`
- `POST /api/redis-management/keys/delete`
- `POST /api/redis-management/keys/expire`
- `POST /api/redis-management/keys/db-stats`
- `POST /api/redis-management/cluster-check`

网络与工具：

- `POST /api/ping/test`
- `POST /api/telnet/test`
- `POST /api/trace-route/test`
- `POST /api/port-scan/scan`

AI 与历史记录：

- `POST /api/ai/chat`
- `POST /api/ai/test-connection`
- `GET /api/md5-histories`
- `POST /api/md5-histories`
- `DELETE /api/md5-histories/:id`
- `DELETE /api/md5-histories`
- `GET /api/image-compress-histories`
- `POST /api/image-compress-histories`
- `DELETE /api/image-compress-histories/:id`
- `DELETE /api/image-compress-histories`

Jenkins 代理：

- `ANY /api/jenkins/proxy/*path`

---

## 10. 宣传文案建议

### 10.1 一句话介绍

DeanTech 是一套面向企业基础设施团队的一体化运维管理平台，集 Kubernetes 管理、主机与虚拟化、容器运维、告警中心、日志查询、协同编辑与批量自动化于一体。

### 10.2 官网短介绍

DeanTech 通过统一的权限体系、统一的数据视图和统一的操作入口，帮助企业构建标准化、可审计、可扩展的运维中台。平台覆盖集群、主机、容器、告警、日志、文档、数据库与自动化执行等核心场景，适用于平台团队、SRE 团队、运维团队和内部技术支持团队。

### 10.3 面向客户的价值表达

- 降低运维工具分散带来的管理成本
- 提升故障排查与日常操作效率
- 建立统一权限和审计机制
- 打通资源管理、自动化执行与协同文档
- 支撑企业内部平台化建设和标准化交付

### 10.4 适合宣传页展示的能力标签

- 多集群 Kubernetes 管理
- 主机与 KVM 统一运维
- Docker 与 Compose 在线管理
- 告警、日志与证书监控
- 批量任务与自动化执行
- 文档协同与图床资产管理
- Redis 运维与常用网络工具
- 权限审计与系统配置中心

---

## 11. 对外发布建议

如果准备正式宣传平台，建议同时准备以下材料：

- 平台首页截图
- 集群总览截图
- 主机 Web SSH / Web RDP 截图
- 图床管理与图片资产页截图
- Redis 管理截图
- 批量运维任务截图
- 文档协同截图
- 一页式卖点海报
- 版本说明与更新日志

建议将本文档作为“完整版产品能力与部署手册”，再单独拆出以下对外材料：

- 一页式产品介绍
- 销售宣讲 PPT
- 标准部署手册
- API 对接手册
- FAQ 与常见问题清单
