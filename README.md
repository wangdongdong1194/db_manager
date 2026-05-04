# use_go

使用 Go + Gin 构建的静态资源内嵌服务示例，web 目录会被编译进二进制文件。

## 目录结构

- main.go: 程序入口，负责 embed 与启动
- internal/config/config.go: 环境变量加载与配置初始化
- internal/router/router.go: 路由注册
- internal/handler/health.go: 示例 API 处理器
- web/: 前端静态资源
- .env.example: 环境变量模板

## 环境准备

1. 安装依赖

go mod tidy

2. 创建本地环境变量文件

cp .env.example .env

3. 按需修改 .env

DB_HOST=localhost
SERVER_PORT=8080
GIN_MODE=release
TRUSTED_PROXIES=127.0.0.1,::1

## 开发运行

go run .

## 路由示例

- /: 前端首页（embed 静态文件）
- /api/health: 健康检查接口

## 打包运行

go build -o main .
./main
