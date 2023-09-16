FROM ubuntu:22.04

ENV GO_VERSION 1.21.0

RUN set -eux; \
       apt-get update; \
       apt-get install -y --no-install-recommends \
               g++ \
               gcc \
               libc6-dev \
               make \
               pkg-config \
               ca-certificates \
               wget \
               curl \
               vim


RUN apt-get install -y git 

RUN rm -rf /var/lib/apt/lists/* 

RUN cd /opt/; \
       curl -OL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
       sha256sum go${GO_VERSION}.linux-amd64.tar.gz && \
       tar -C /usr/local -xvf go${GO_VERSION}.linux-amd64.tar.gz && \
       rm -rf /opt/go${GO_VERSION}.linux-amd64.tar.gz

 
ENV PATH /usr/local/go/bin:$PATH

CMD ["go", "run", "./chat/web-socket-chat/main.go"]

WORKDIR / 

