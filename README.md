# golang-common-base Monorepo

当前仓库已重构为三层结构：

- `backend-admin`：Go 后端服务（Gin + Gorm）
- `admin-frontend`：管理端前端（基于 v3-admin-vite）
- `deploy`：统一部署目录（本地/容器编排）

## 1. 后端

后端目录：`backend-admin`

关键能力：
- 数据库：MySQL / SQLite（配置切换）
- 缓存：Redis
- 对象存储：MinIO
- 身份：Casdoor（OIDC）
- 授权：Casbin（RBAC）

快速编译：

```bash
cd backend-admin
go test ./...
```

## 2. 前端

前端目录：`admin-frontend`

来源：
- https://github.com/un-pany/v3-admin-vite

启动方式请参考 `admin-frontend` 目录下文档。

## 3. 部署

部署目录：`deploy`

本地容器化启动（MySQL + Redis + MinIO + Casdoor + Backend）：

```bash
cd deploy
docker compose up -d --build
```

详细说明见：`deploy/README.md`
