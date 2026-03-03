# DeanTech 系统文档

## 1. 系统介绍

DeanTech是一款综合性的监控管理平台，提供了从基础设施到应用层的全方位监控和管理功能。该系统支持Kubernetes集群管理、主机监控、告警管理、日志分析等核心功能，旨在帮助用户实现对IT基础设施的集中监控和高效管理。

### 1.1 主要功能

- **集群管理**：Kubernetes集群节点、Pod、Deployment、Service等资源的监控和管理
- **主机监控**：支持物理机和虚拟机的监控
- **告警管理**：集成Prometheus告警，支持告警规则配置和告警历史查询
- **日志管理**：支持日志查询和分析
- **邮箱配置**：系统告警邮件发送配置
- **系统设置**：全局系统参数配置
- **用户管理**：支持多用户系统，包含登录、注册、密码重置等功能

## 2. 技术栈

### 2.1 前端技术栈

| 技术/框架 | 版本 | 用途 |
|---------|-----|-----|
| Vue.js | ^3.4.21 | 前端框架 |
| Vue Router | ^4.3.0 | 路由管理 |
| TypeScript | ^5.2.2 | 类型系统 |
| Vite | ^5.2.0 | 构建工具 |
| Element Plus | ^2.6.1 | UI组件库 |
| Axios | ^1.6.7 | HTTP客户端 |
| ECharts | ^5.5.0 | 图表库 |
| Xterm.js | ^5.3.0 | 终端模拟器 |

### 2.2 后端技术栈

| 技术/框架 | 版本 | 用途 |
|---------|-----|-----|
| Go | 1.23.0 | 后端开发语言 |
| Gin | v1.9.1 | Web框架 |
| GORM | v1.30.0 | ORM框架 |
| MySQL | - | 关系型数据库 |
| SQLite | - | 轻量级数据库（用于Prometheus告警） |
| Kubernetes Client-go | v0.28.0 | Kubernetes API客户端 |
| Prometheus | - | 监控系统集成 |
| WebSocket | - | 实时通信 |
| SFTP | - | 文件传输 |

## 3. 部署指南

### 3.1 环境要求

- **操作系统**：Linux/macOS/Windows
- **Go版本**：1.23.0+（仅后端开发/编译）
- **Node.js版本**：18.x+（仅前端开发/编译）
- **数据库**：MySQL 5.7+ 或 SQLite（默认）
- **容器化**：Docker（可选，用于镜像构建）

### 3.2 前端部署

#### 3.2.1 开发环境

1. 进入前端目录：
   ```bash
   cd /Users/dean/DeanTech/frontend
   ```

2. 安装依赖：
   ```bash
   npm install
   ```

3. 启动开发服务器：
   ```bash
   npm run dev
   ```

4. 访问前端应用：
   ```
   http://localhost:5173
   ```

#### 3.2.2 生产构建

1. 编译打包：
   ```bash
   npm run build
   ```

2. 构建产物位于 `dist/` 目录，可直接部署到Web服务器。

#### 3.2.3 部署到Nginx

示例Nginx配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/deantech/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 3.3 后端部署

#### 3.3.1 开发环境

1. 进入后端目录：
   ```bash
   cd /Users/dean/DeanTech/backend
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 启动后端服务：
   ```bash
   go run cmd/main.go
   ```

4. 后端服务默认运行在 `http://localhost:8080`。

#### 3.3.2 编译构建

1. 编译可执行文件：
   ```bash
   # Linux
   GOOS=linux GOARCH=amd64 go build -o deantech-backend cmd/main.go
   
   # macOS
   GOOS=darwin GOARCH=amd64 go build -o deantech-backend cmd/main.go
   
   # Windows
   GOOS=windows GOARCH=amd64 go build -o deantech-backend.exe cmd/main.go
   ```

2. 运行编译后的可执行文件：
   ```bash
   ./deantech-backend
   ```

#### 3.3.3 Docker镜像构建

1. 创建Dockerfile（后端）：

```dockerfile
# 后端Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o deantech-backend cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/deantech-backend .
EXPOSE 8080
CMD ["./deantech-backend"]
```

2. 构建后端镜像：
   ```bash
   docker build -t deantech-backend:latest .
   ```

3. 创建Dockerfile（前端）：

```dockerfile
# 前端Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

4. 构建前端镜像：
   ```bash
   cd /Users/dean/DeanTech/frontend
   docker build -t deantech-frontend:latest .
   ```

5. 使用Docker Compose部署：

创建 `docker-compose.yml` 文件：

```yaml
version: '3'
services:
  backend:
    image: deantech-backend:latest
    ports:
      - "8080:8080"
    volumes:
      - ../data:/app/data
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=deantech
    depends_on:
      - mysql

  frontend:
    image: deantech-frontend:latest
    ports:
      - "80:80"
    depends_on:
      - backend

  mysql:
    image: mysql:5.7
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=deantech

volumes:
  mysql-data:
```

启动服务：
```bash
docker-compose up -d
```

### 3.4 数据库初始化

系统默认使用SQLite数据库，位于 `/Users/dean/DeanTech/data/promalert.db`，会在首次运行时自动创建。

如果需要使用MySQL数据库，系统会在首次运行时自动创建所需的表结构。以下是手动创建数据库和表的SQL语句，可供参考：

#### 3.4.1 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS deantech CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 3.4.2 用户表

```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    department VARCHAR(50) NULL,
    position VARCHAR(50) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 3.4.3 邮箱设置表

```sql
CREATE TABLE IF NOT EXISTS email_settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    smtp_server VARCHAR(100) NOT NULL,
    smtp_port INT NOT NULL,
    from_email VARCHAR(100) NOT NULL,
    smtp_username VARCHAR(100) NOT NULL,
    smtp_password VARCHAR(255) NOT NULL,
    smtp_encryption VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 3.4.4 系统设置表

```sql
CREATE TABLE IF NOT EXISTS system_settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    key_name VARCHAR(50) NOT NULL UNIQUE,
    key_value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 4. 系统对接使用

### 4.1 初始登录

系统首次运行时，默认会创建一个管理员账户：

- **用户名**：admin
- **密码**：admin123

首次登录后建议立即修改密码。

### 4.2 主要功能使用

#### 4.2.1 集群管理

1. 登录系统后，点击左侧菜单的"集群管理"
2. 可以查看集群节点、Pod、Deployment、Service等资源
3. 支持查看资源详情和基本操作

#### 4.2.2 告警管理

1. 点击左侧菜单的"告警管理"
2. 可以查看当前告警和告警历史
3. 支持配置告警规则

#### 4.2.3 邮箱配置

1. 点击左侧菜单的"系统设置" -> "邮箱配置"
2. 填写SMTP服务器信息
3. 点击"测试发送邮件"验证配置是否正确
4. 点击"保存配置"保存设置

#### 4.2.4 用户管理

1. 点击左侧菜单的"系统设置" -> "用户管理"
2. 可以添加、编辑、删除用户
3. 支持修改用户密码

### 4.3 API接口

系统提供RESTful API接口，前缀为 `/api`，主要接口包括：

- **用户管理**：`/api/users/*`
- **集群管理**：`/api/cluster/*`
- **告警管理**：`/api/alerts/*`
- **邮箱配置**：`/api/email/*`
- **系统设置**：`/api/system/*`

具体接口文档可通过访问系统的Swagger页面查看（如果已配置）。

## 5. 常见问题

### 5.1 如何修改默认端口？

- **前端**：修改 `vite.config.ts` 中的 `server.port` 配置
- **后端**：修改 `cmd/main.go` 中的端口配置，默认8080

### 5.2 如何切换数据库？

修改后端配置文件 `config/config.go` 中的数据库连接信息。

### 5.3 如何配置HTTPS？

建议通过Nginx或其他反向代理服务器配置HTTPS，后端服务保持HTTP即可。

## 6. 维护与更新

### 6.1 日志查看

- **前端**：浏览器开发者工具的控制台
- **后端**：默认输出到控制台，可配置日志文件

### 6.2 升级步骤

1. 备份数据（特别是数据库文件）
2. 停止当前运行的服务
3. 更新代码或替换可执行文件
4. 启动服务
5. 验证功能是否正常

## 7. 联系方式

如有问题或建议，欢迎联系系统管理员或开发团队。

---

**文档版本**：v1.0
**更新日期**：2026-02-14
**编写人员**：DeanTech开发团队


curl http://localhost:8000/api/dashboard/health-status
{"systemServices":{"percentage":87,"status":"warning","description":"核心服务运行状态"},"apiResponse":{"percentage":89,"status":"good","description":"API请求响应时间"},"databaseConnection":{"percentage":97,"status":"good","description":"数据库连接稳定性"}}%