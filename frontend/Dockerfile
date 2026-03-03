# 构建阶段
FROM docker.m.daocloud.io/node:18-alpine as builder
WORKDIR /app

# 复制依赖文件
COPY package*.json ./

# 安装依赖
#RUN npm ci
RUN npm install

# 复制源代码
COPY . .

# 构建生产版本
RUN npm run build

# 运行阶段
FROM docker.m.daocloud.io/nginx:alpine

# 复制构建好的静态文件到Nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制自定义Nginx配置（可选）
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]
