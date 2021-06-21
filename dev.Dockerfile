FROM golang:1.15
MAINTAINER Mustang <mustang2247@gmail.com>

ENV PORT_TO_EXPOSE=9088
ENV PROC_NAME=golang-common-base
ENV ENV_SERVER_MODE=dev

#设置工作目录
WORKDIR /opt/app
# 将服务器的go工程代码加入到docker容器中
COPY . .

#RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN go mod download
# go构建可执行文件
#RUN go build .

COPY docker/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >  /etc/timezone

#暴露端口
EXPOSE 9088
ENTRYPOINT ["air", "-d"]
#ENTRYPOINT ["./golang-common-base"]
