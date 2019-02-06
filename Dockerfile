# vi: set ft=dockerfile :

FROM centos:7
MAINTAINER Kazuhisa Hara <khara@sios.com>

RUN yum install -y --setopt=tsflags=nodocs \
        epel-release \
        bash bash-completion && \
    yum install -y --setopt=tsflags=nodocs \
        golang && \
    rm -rf /var/cache/yum

WORKDIR /root

RUN go get -u github.com/go-oauth2/gin-server

CMD /bin/bash
