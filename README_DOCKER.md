# DeanTech Docker 部署说明

## 项目结构

```
DeanTech/
├── backend/          # 后端 Go 应用
│   ├── Dockerfile    # 后端 Dockerfile
├── frontend/         # 前端 Vue 应用
│   ├── Dockerfile    # 前端 Dockerfile
└── docker-compose.yml # Docker Compose 配置文件
```

## 构建和运行

### 方法一：使用 Docker Compose（推荐）

Docker Compose 可以同时构建和运行前端和后端服务，是最便捷的方式。

1. **构建并启动服务**：

```bash
# 进入项目根目录
cd /Users/dean/DeanTech

# 构建并启动服务
# -d 表示后台运行
# --build 表示强制重新构建镜像
docker-compose up -d --build
```

2. **停止服务**：

```bash
docker-compose down
```

3. **查看日志**：

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs frontend
docker-compose logs backend
```

### 方法二：手动构建和运行

#### 前端构建和运行

1. **构建前端镜像**：

```bash
cd /Users/dean/DeanTech/frontend
docker build -t deantech-frontend .
```

2. **运行前端容器**：

```bash
docker run -d \
  -p 80:80 \
  --name deantech-frontend \
  --network deantech-network \
  -e VUE_APP_API_BASE_URL=http://backend:8000/api \
  deantech-frontend
```

#### 后端构建和运行

1. **构建后端镜像**：

```bash
cd /Users/dean/DeanTech/backend
docker build -t deantech-backend .
```

2. **运行后端容器**：

```bash
docker run -d \
  -p 8000:8000 \
  --name deantech-backend \
  --network deantech-network \
  -e GIN_MODE=release \
  -e DATABASE_HOST=db \
  -e DATABASE_PORT=3306 \
  -e DATABASE_USER=root \
  -e DATABASE_PASSWORD=password \
  -e DATABASE_NAME=deantech \
  deantech-backend
```

## 环境变量配置

### 前端环境变量

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| VUE_APP_API_BASE_URL | 后端 API 基础 URL | http://backend:8000/api |

### 后端环境变量

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| GIN_MODE | Gin 运行模式（debug/release/test） | release |
| DATABASE_HOST | 数据库主机地址 | db |
| DATABASE_PORT | 数据库端口 | 3306 |
| DATABASE_USER | 数据库用户名 | root |
| DATABASE_PASSWORD | 数据库密码 | password |
| DATABASE_NAME | 数据库名称 | deantech |

## 自定义配置

### 配置文件挂载

如果需要使用自定义配置文件，可以通过挂载方式替换容器内的配置：

```bash
# 后端配置文件挂载示例
docker run -d \
  -p 8000:8000 \
  --name deantech-backend \
  -v /path/to/your/config.yaml:/app/config.yaml \
  deantech-backend
```

### 端口映射

如果需要使用不同的端口，可以修改端口映射：

```bash
# 前端使用 8080 端口
docker run -d -p 8080:80 deantech-frontend

# 后端使用 9000 端口
docker run -d -p 9000:8000 deantech-backend
```

## 常见问题

1. **前端无法连接到后端**
   - 确保前后端在同一个网络中
   - 检查环境变量 `VUE_APP_API_BASE_URL` 是否正确
   - 检查后端服务是否正常运行

2. **构建失败**
   - 确保 Docker 版本足够新
   - 检查网络连接，确保可以下载依赖
   - 检查代码是否有语法错误

3. **端口被占用**
   - 使用 `docker ps` 查看是否有其他容器占用了相同端口
   - 修改端口映射，使用其他端口

## 更新镜像

当代码更新后，需要重新构建镜像：

```bash
# 使用 Docker Compose 更新
docker-compose up -d --build

# 手动更新前端镜像
docker build -t deantech-frontend .
docker run -d -p 80:80 deantech-frontend

# 手动更新后端镜像
docker build -t deantech-backend .
docker run -d -p 8000:8000 deantech-backend
```

## 数据持久化

如果需要持久化数据，比如数据库数据，可以在 `docker-compose.yml` 中添加数据库服务和数据卷配置。

## 更多信息

- Docker 官方文档：https://docs.docker.com/
- Docker Compose 官方文档：https://docs.docker.com/compose/
