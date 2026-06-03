# 部署说明

本项目推荐优先使用 `docker compose` 部署。后端依赖 MySQL，前端通过 Nginx 反向代理 `/api` 到后端，`guacd` 和 `guac-gateway` 用于主机/控制台相关能力。

## 1. 部署前准备

- Docker 24+
- Docker Compose v2
- MySQL 8.0+
- 如需 K8s 管理能力，目标集群需可访问
- 如需 Prometheus 页面，`prometheus.addr` 必须可访问

## 2. Docker 方式

### 2.1 运行

```bash
docker network create deantech-network

docker run -d \
  --name mysql \
  --network deantech-network \
  --network-alias mysql \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -e MYSQL_DATABASE=deantech \
  -e MYSQL_USER=deantech \
  -e MYSQL_PASSWORD=deantech123 \
  -v mysql_data:/var/lib/mysql \
  mysql:8.0

docker run -d \
  --name backend \
  --network deantech-network \
  -p 8000:8000 \
  -e GIN_MODE=release \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8000 \
  -e DATABASE_DSN='deantech:deantech123@tcp(mysql:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local' \
  -e PROMETHEUS_ADDRESS='http://prometheus:9090' \
  -e ALERTMANAGER_URL='http://alertmanager:9093' \
  registry.cn-beijing.aliyuncs.com/deanmr/deantech:deantech-backend-v6.0

> 推荐把 `kubectl` 作为平台发布资产一起构建进镜像，而不是依赖容器基础镜像或宿主机环境。

> 如果你的 Prometheus / Alertmanager 不在同一个 Docker 网络中，请把上述地址替换成实际可访问的地址。

docker run -d \
  --name guacd \
  --network deantech-network \
  registry.cn-beijing.aliyuncs.com/deanmr/deantech:guacd

docker run -d \
  --name guac-gateway \
  --network deantech-network \
  -p 8081:8081 \
  -e PORT=8081 \
  -e GUACD_HOST=guacd \
  -e GUACD_PORT=4822 \
  -e GUACAMOLE_TOKEN_KEY=deantech-guacamole-secret \
  registry.cn-beijing.aliyuncs.com/deanmr/deantech:guac-gateway

docker run -d \
  --name frontend \
  --network deantech-network \
  -p 80:80 \
  registry.cn-beijing.aliyuncs.com/deanmr/deantech:deantech-frontend-v6.0
```

## 3. docker-compose 方式

仓库根目录的 [docker-compose.yml](docker-compose.yml) 已经补齐 MySQL、backend、frontend、guacd 和 guac-gateway。
Prometheus 和 Alertmanager 默认按外部服务处理，请把 `PROMETHEUS_ADDRESS` 和 `ALERTMANAGER_URL` 指向你真实可访问的地址。

### 3.1 启动

```bash
docker compose up -d
```


### 3.2 停止

```bash
docker compose down
```

### 3.3 查看日志

```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
```

### 3.4 服务说明

| 服务 | 说明 | 端口 |
| --- | --- | --- |
| `mysql` | 主数据库 | `3306` |
| `backend` | Go API 服务 | `8000` |
| `frontend` | Vue + Nginx | `80` |
| `guacd` | Guacamole 代理 | `4822` |
| `guac-gateway` | WebSocket 网关 | `8081` |

### 3.5 关键环境变量

| 环境变量 | 说明 |
| --- | --- |
| `DATABASE_DSN` | MySQL 连接串 |
| `PROMETHEUS_ADDRESS` | Prometheus 地址 |
| `ALERTMANAGER_URL` | Alertmanager 地址 |
| `GUACAMOLE_WS_URL` | 网关 WebSocket 地址 |
| `GUACAMOLE_TOKEN_KEY` | 网关令牌密钥 |
| `DEANTECH_KUBECTL_BIN` | 可选：显式指定平台内置 kubectl 路径 |
| `GUACAMOLE_TUNNEL_HOST` | `guacd` 访问临时 RDP 转发口时使用的主机名 |
| `GUACAMOLE_TUNNEL_BIND` | 后端临时 RDP 转发监听地址 |
| `GUACAMOLE_TUNNEL_TTL_SECONDS` | 临时 RDP 转发保留时长，默认 `600` 秒 |

### 3.6 后端离线发布包

后端支持打成自包含发布目录，目录内会保留 `config/`、`migrations/` 和平台内置 `kubectl` 目录：

```bash
cd backend
TARGET_OS=linux TARGET_ARCH=amd64 VERSION=v1.3.0 ./scripts/package_release.sh
```

输出：

- `dist/deantech-backend-v1.3.0-linux-amd64/`
- `dist/deantech-backend-v1.3.0-linux-amd64.tar.gz`

## 4. Guacamole 控制台服务说明

本项目的 Windows 控制台链路由 `backend -> guac-gateway -> guacd -> Windows RDP` 组成：

- `backend` 负责根据主机信息生成 Guacamole 加密 token。
- `guac-gateway` 提供浏览器 WebSocket 接入，并把 token 解密后转发给 `guacd`。
- `guacd` 负责实际发起 RDP 连接并把画面编码为 Guacamole 协议。
- `frontend` 通过 `/guac-ws/` 连接 `guac-gateway`，在页面内嵌 Web RDP 控制台。

### 4.1 部署要求

- `guacd` 与 `guac-gateway` 必须能互相访问。
- `backend` 与 `guac-gateway` 必须使用相同的 `GUACAMOLE_TOKEN_KEY`。
- 前端代理必须转发 `/guac-ws/` 到 `guac-gateway:8081`。
- 目标 Windows 主机必须开启 RDP，且防火墙允许 `3389` 或你自定义的 RDP 端口。
- 如果 Windows 虚机需要经“网域主机”跳转，`guacd` 必须能访问后端暴露出来的临时转发地址。

### 4.2 docker compose 推荐方式

仓库根目录的 [docker-compose.yml](docker-compose.yml) 已经包含：

- `guacd`
- `guac-gateway`
- `backend` 中的 `GUACAMOLE_WS_URL=/guac-ws/`
- `backend` / `guac-gateway` 共享的 `GUACAMOLE_TOKEN_KEY`
- KVM Windows 虚机经网域主机跳转时，容器内默认会使用 `backend` 服务名访问临时转发口

启动命令：

```bash
docker compose up -d --build guacd guac-gateway backend frontend
```

检查状态：

```bash
docker compose ps guacd guac-gateway
docker compose logs -f guacd
docker compose logs -f guac-gateway
```

### 4.3 路径与端口约定

| 组件 | 默认值 | 说明 |
| --- | --- | --- |
| `guacd` | `4822/tcp` | Guacamole 代理守护进程 |
| `guac-gateway` | `8081/tcp` | 浏览器接入网关 |
| `GUACAMOLE_WS_URL` | `/guac-ws/` | 前端访问路径 |
| 前端 Vite 代理 | `127.0.0.1:8081` | 本地开发代理 |
| 前端 Nginx 代理 | `guac-gateway:8081` | 容器部署代理 |

### 4.4 Windows 主机字段怎么填

如果你想使用系统内置的 Web RDP 控制台，Windows 主机应这样填写：

- `IP地址`：Windows 主机真实地址
- `端口`：通常为 `3389`
- `用户名`：Windows 登录用户名
- `密码`：Windows 登录密码
- `控制台地址`：留空

说明：

- 当 `控制台地址` 留空时，系统会自动走内置 Guacamole 模式。
- 当 `控制台地址` 填写为 `http://` 或 `https://` 网页地址时，系统会改为外部 `iframe` 嵌入模式。
- 不要填写 `rdp://...`、`/guac-ws/`、`guacd:4822` 这类值。

### 4.5 验证步骤

1. 确认 `guac-gateway` 已监听 `8081`。
2. 确认 `guac-gateway` 容器内部可以连接到 `guacd:4822`。
3. 在前端打开“虚机管理 -> Windows 主机 -> 控制台”。
4. 页面应显示“已连接”，并能看到 Windows 桌面画面。

可用验证命令：

```bash
docker compose exec -T guac-gateway node -e "const net=require('net');const s=net.createConnection({host:'guacd',port:4822,timeout:3000},()=>{console.log('CONNECTED');s.end()});s.on('error',e=>{console.error(e.message);process.exit(1)})"
```

### 4.6 常见问题排查

#### 浏览器显示“未连接”

- 检查 `backend` 是否正确返回 `/guac-ws/` 和 token。
- 检查前端代理是否已把 `/guac-ws/` 转发到 `guac-gateway`。
- 检查 `backend` 与 `guac-gateway` 的 `GUACAMOLE_TOKEN_KEY` 是否一致。

#### `guac-gateway` 启动正常，但无法连接 `guacd`

- 如果你是在宿主机直接运行 `npm start`，默认 `GUACD_HOST=guacd` 可能无法解析。
- 推荐直接使用 `docker compose` 启动 `guacd + guac-gateway`。
- 若坚持宿主机运行 `guac-gateway`，需要显式设置可达的 `GUACD_HOST`，例如 `127.0.0.1`。

#### KVM Windows 虚机启用了网域主机，但浏览器里仍然连不上

- Docker Compose 部署下，后端会默认使用 `backend` 作为临时 RDP 转发主机名。
- 如果你是“后端跑宿主机、guacd/guac-gateway 跑容器”的混合部署，通常需要在后端显式设置：
  - `GUACAMOLE_TUNNEL_HOST=host.docker.internal`
  - `GUACAMOLE_TUNNEL_BIND=0.0.0.0`
- 如果 `guacd` 与后端不在同一网络，请把 `GUACAMOLE_TUNNEL_HOST` 改成 `guacd` 实际可达的后端地址。

#### Windows 控制台显示已连接，但黑屏

- 优先确认目标 Windows 是否允许 RDP 登录，且没有被其他会话强制顶掉。
- 检查 `guacd` 日志中是否出现 `ERRINFO_DISCONNECTED_BY_OTHER_CONNECTION`。
- 如果页面已有鼠标但没有画面，通常需要同时检查：
  - RDP 协议参数兼容性
  - 浏览器端 Guacamole display 缩放/适配
  - Windows 目标会话是否真实输出桌面
- 当前项目已启用一套更保守的兼容参数与 display 自适配逻辑；如仍异常，优先点击“重连”。

#### `控制台地址` 填了值但打不开

- 该字段只支持可直接嵌入的网页地址。
- 如果目标站点设置了 `X-Frame-Options` 或 CSP 禁止嵌入，页面会无法显示。
- 这类场景建议清空 `控制台地址`，改用内置 Web RDP。

## 5. Kubernetes 方式

### 5.1 推荐拓扑

- `Namespace`
- `Secret`：数据库密码、JWT、Guacamole token
- `ConfigMap`：后端配置、前端 Nginx 配置
- `StatefulSet` 或外部 MySQL
- `Deployment`：backend、frontend、guacd、guac-gateway
- `Service`：backend、frontend、mysql、guacd、guac-gateway
- `Ingress`：统一对外入口

### 5.2 部署顺序

1. 创建命名空间和密钥。
2. 部署 MySQL。
3. 部署 backend。
4. 部署 guacd 和 guac-gateway。
5. 部署 frontend。
6. 配置 Ingress。

### 5.3 示例 YAML

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: deantech
---
apiVersion: v1
kind: Secret
metadata:
  name: deantech-secret
  namespace: deantech
type: Opaque
stringData:
  mysql-root-password: root123
  mysql-password: deantech123
  guacamole-token-key: deantech-guacamole-secret
  jwt-secret: your-secret-key
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: deantech-config
  namespace: deantech
data:
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "8000"
  DATABASE_DSN: "deantech:deantech123@tcp(mysql:3306)/deantech?charset=utf8mb4&parseTime=True&loc=Local"
  PROMETHEUS_ADDRESS: "http://prometheus.monitoring.svc.cluster.local:9090"
  ALERTMANAGER_URL: "http://alertmanager.monitoring.svc.cluster.local:9093"
```

### 5.4 backend Deployment 示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-backend
  namespace: deantech
spec:
  replicas: 1
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
          image: registry.cn-beijing.aliyuncs.com/deanmr/deantech:deantech-backend-v6.0
          ports:
            - containerPort: 8000
          envFrom:
            - configMapRef:
                name: deantech-config
          env:
            - name: GIN_MODE
              value: release
          readinessProbe:
            httpGet:
              path: /health
              port: 8000
            initialDelaySeconds: 10
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /health
              port: 8000
            initialDelaySeconds: 30
            periodSeconds: 20
```

### 5.5 frontend Deployment 示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deantech-frontend
  namespace: deantech
spec:
  replicas: 1
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
          image: registry.cn-beijing.aliyuncs.com/deanmr/deantech:deantech-frontend-v6.0
          ports:
            - containerPort: 80
```

### 5.6 Ingress 示例

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: deantech
  namespace: deantech
spec:
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
```

## 6. 验证

- 打开首页后确认登录页/首页配置是否生效。
- 打开集群总览确认能拉到集群列表与节点数据。
- 打开 Prometheus 页确认 `prometheus.addr` 可访问。
- 打开主机管理确认 WebSocket 和 `guac-gateway` 可正常工作。
