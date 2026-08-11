# Deploy 目录说明

## 目录结构

- `docker-compose.yaml`：本地一键启动后端与基础依赖（MySQL、Redis、MinIO、Casdoor）。
- `docker-compose.casdoor-prod.yaml`：生产模式 Casdoor 独立部署编排（参考官方 Docker 部署）。
- `caddy/Caddyfile`：Caddy 反向代理配置（自动 HTTPS）。
- `casdoor/README.md`：Casdoor 生产部署与使用手册。
- `.env.casdoor.example`：Casdoor 生产部署环境变量示例。

## 快速启动

1. 进入 deploy 目录。
2. 执行：

```bash
docker compose up -d --build
```

3. 访问服务：

- Backend API: http://127.0.0.1:9080
- MinIO API: http://127.0.0.1:9000
- MinIO Console: http://127.0.0.1:9001
- Casdoor: http://127.0.0.1:8000

## Casdoor 生产部署

参考官方文档：

- https://casdoor.ai/zh/docs/deployment/docker

本仓库提供了生产化模板，包含：

- 独立 `casdoor` 服务，数据保存到 `.env.casdoor` 指定的外部 MySQL（默认库名 `db_casdoor`）
- `GIN_MODE=release`
- 配置与日志卷映射（`./casdoor/conf:/conf`, `./casdoor/logs:/logs`）
- 环境变量外置（`.env.casdoor`）
- Caddy 自动签发与续期 HTTPS 证书

启动命令：

```bash
cd deploy
cp .env.casdoor.example .env.casdoor
cp casdoor/conf/app.conf.example casdoor/conf/app.conf
docker compose --env-file .env.casdoor -f docker-compose.casdoor-prod.yaml up -d
```

使用前请确认：
- `.env.casdoor` 中 `CASDOOR_DOMAIN` 已改为真实域名
- 域名 DNS 已指向部署机器
- 80/443 端口已放通

详细步骤见：`casdoor/README.md`

## 使用建议

- 开发联调：使用 `docker-compose.yaml`。
- 生产环境：优先使用 `docker-compose.casdoor-prod.yaml`，并在反向代理层配置 HTTPS。
- 该生产模板默认集成 Caddy，自动支持 HTTPS。
- 后端对接时，`auth.casdoor.issuer` 应填写 Casdoor 对外 HTTPS 域名，不能使用内网地址。

## 注意

- 首次启动前需要创建数据库 `db_casdoor`，并确保 `.env.casdoor` 中的账号具有访问权限。
- 如果你使用 SQLite，请在 `backend-admin/conf/config_dev.yaml` 将 `db.driver` 改为 `sqlite` 并设置 `db.sqlite.path`。
- Casdoor 作为身份平台，Casbin 用于后端授权。你需要在 Casdoor 中创建应用并配置 `client_id` 与 issuer。
