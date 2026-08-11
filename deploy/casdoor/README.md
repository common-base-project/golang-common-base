# Casdoor 生产部署（Docker）

本文档基于 Casdoor 官方 Docker 部署说明整理：
- https://casdoor.ai/zh/docs/deployment/docker

## 1. 前置条件

- Docker
- Docker Compose
- 生产域名（建议）
- 80/443 端口（若使用反向代理 + TLS）

如果使用本文提供的 Caddy 自动 HTTPS 方案，必须满足：
- `CASDOOR_DOMAIN` 的 DNS A/AAAA 记录指向当前服务器
- 防火墙/安全组放通 80 和 443

## 2. 准备配置目录

当前仓库已提供：
- `deploy/casdoor/conf`
- `deploy/casdoor/logs`

建议将官方配置模板拉取到 conf 目录：

```bash
cd deploy
cp casdoor/conf/app.conf.example casdoor/conf/app.conf
# 可选：拉取官方模板比对
# wget https://raw.githubusercontent.com/casdoor/casdoor/master/conf/app.conf -O casdoor/conf/app.conf
# wget https://raw.githubusercontent.com/casdoor/casdoor/master/init_data.json.template -O casdoor/conf/init_data.json
```

## 3. 配置数据库连接

1. 复制环境变量模板：

```bash
cd deploy
cp .env.casdoor.example .env.casdoor
```

2. 确认 `.env.casdoor` 中的 MySQL 地址、端口和账号，并填写密码。默认指向 `backend-admin/conf/config_dev.yaml` 使用的 MySQL 服务，数据库名为 `db_casdoor`。
3. 修改 `.env.casdoor` 中 `CASDOOR_DOMAIN` 与 `CADDY_ACME_EMAIL`。

Casdoor 的 `dataSourceName` 由 Compose 根据 `.env.casdoor` 生成，无需再手动修改 `casdoor/conf/app.conf`。启动前请确保 `db_casdoor` 已创建，且该账号对其有权限。

建议在 `casdoor/conf/app.conf` 中设置 Casdoor 对外地址（例如 `https://auth.example.com`），确保与 OIDC issuer 一致。

## 4. 启动 Casdoor（生产模式）

```bash
cd deploy
docker compose --env-file .env.casdoor -f docker-compose.casdoor-prod.yaml up -d
```

该命令会同时启动：
- `casdoor`（身份服务）
- `casdoor-caddy`（反向代理 + 自动申请/续期 HTTPS 证书）

查看日志：

```bash
docker compose --env-file .env.casdoor -f docker-compose.casdoor-prod.yaml logs -f casdoor
docker compose --env-file .env.casdoor -f docker-compose.casdoor-prod.yaml logs -f caddy
```

## 5. 访问与联通性验证

默认通过 Caddy 对外提供 HTTPS：
- https://你的域名

基础检查：

```bash
docker ps
curl -I https://你的域名
```

## 6. 与 backend-admin 对接

在后端配置中设置：
- `auth.casdoor.enabled: true`
- `auth.casdoor.issuer: https://你的域名`
- `auth.casdoor.client_id: 你的应用ID`

注意：`issuer` 必须与 Casdoor 实际对外地址一致（建议 HTTPS 域名）。

## 7. 常见问题

- 容器启动失败：先看 Casdoor 与 MySQL 日志。
- 登录回调失败：检查 Casdoor 应用回调地址是否和前端/后端一致。
- Token 验证失败：通常是 `issuer` 或 `client_id` 不匹配。
- 文件权限问题：Casdoor 容器默认 UID/GID 为 1000，映射目录需有写权限。
- HTTPS 证书申请失败：通常是域名未解析到当前机器，或 80/443 端口被占用/未放通。
