# Golang + Gin + Gorm + MySQL 基础手脚架
    Golang + Gin + Gorm + MySQL 基础手脚架，支持容器和k8s部署。打包后自动生成 docker image 可以直接配置好账户，直接自动化上传。

## make 打包
    注意：Makefile 文件里 main.go 的路径

    make docker-all VERSION="staging_v0.0.2" ENV_SERVER_MODE="staging"
    make docker-all VERSION="staging_v0.0.1" ENV_SERVER_MODE="dev"
    make docker-all VERSION="prod_v0.0.1" ENV_SERVER_MODE="prod"

    golang build:
    go build -o golang-common-base ./app/cmd/

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
