FROM alpine
MAINTAINER Mustang <mustang2247@gmail.com>

ARG procname
ARG packagefile
ARG exposeport
ARG modeenv

ENV PORT_TO_EXPOSE=${exposeport}
ENV PROC_NAME=${procname}
ENV ENV_SERVER_MODE=${modeenv}

WORKDIR /opt/${procname}
ADD $packagefile /opt/${procname}/
RUN mkdir -p /opt/${procname}/conf
VOLUME ["/opt/$PROC_NAME"]

COPY Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >  /etc/timezone

CMD ["sh", "-c", "./$PROC_NAME"]
EXPOSE $PORT_TO_EXPOSE

