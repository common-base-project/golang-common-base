# Backend Admin

Golang + Gin + Gorm 基础后端，已支持如下能力：
- 数据库：MySQL、SQLite（可切换）
- 缓存：Redis
- 对象存储：MinIO
- 身份管理：Casdoor（OIDC）
- 权限控制：Casbin（RBAC）

项目已迁移到 monorepo 结构中的 `backend-admin` 目录。

## make 打包
    注意：Makefile 文件里 main.go 的路径

    make docker-all VERSION="staging_v0.0.2" ENV_SERVER_MODE="staging"
    make docker-all VERSION="staging_v0.0.1" ENV_SERVER_MODE="dev"
    make docker-all VERSION="prod_v0.0.1" ENV_SERVER_MODE="prod"

    golang build:
    go build -o golang-common-base ./app/cmd/

## 配置说明

在 `conf/config_dev.yaml` 中可配置：
- `db.driver`: `mysql` 或 `sqlite`
- `db.sqlite.path`: SQLite 数据库路径
- `redis.*`: Redis 连接信息
- `minio.*`: MinIO 连接信息
- `auth.casdoor.*`: Casdoor OIDC 配置
- `authz.casbin.*`: Casbin 模型与策略配置

Casbin 模型文件位置：`conf/casbin_model.conf`。

## 生成`swagger`文档
```
    go get -u github.com/swaggo/swag/cmd/swag
    swag init
# 基于Makefile
    make swagger

# OR 使用 swag 命令（注意：main.go 的路径）
    swag init -g ./app/cmd/main.go  -o ./docs/

```

## 基于 docker 容器开发
```text

# 本项目本地开发步骤：
前提（可选）：
    安装 air 工具: https://github.com/cosmtrek/air

一 直接下载源代码到本地用 IDE 本地调试开发

二 基于 docker 环境开发
    1 安装 docker
    2 下载开发镜像 'golang-common-base:dev_v1' 或者基于源代码编译 docker 镜像
        docker build -f dev.Dockerfile -t golang-common-base:dev_v1 .
        docker push golang-common-base:dev_v1
    
    3 推荐 vscode 基于 docker 开发
        a vscode 需要安装 "Remote - Containers" 工具
        b 这 vscode 编辑器选择快捷键 Cmd + shift + p 输入 "Remote-Containers: Attach to Running Container……" 然后选择 golang-common-base
        c vscode 打开文件夹，打开 "/opt/app" 目录即可开发
        d 可以这 Container 里直接使用自己的git， 也可以直接调试等
        
     
```

## 解决 Mac pro m1 standard_init_linux.go:228: exec user process caused: exec format error
```shell
# 解决 Mac pro m1 （arm芯片）电脑 docker build 默认build是 linux/arm 我们需要 linux/amd64
# 参考文档： https://docs.docker.com/desktop/multi-arch/
docker buildx build --platform linux/amd64
```
