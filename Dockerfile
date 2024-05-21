FROM node:18.11.0 AS builder
ENV NODE_OPTIONS=--openssl-legacy-provider
LABEL stage=build
WORKDIR /build
RUN yarn config set registry https://registry.npmmirror.com/
ENV GOLANG_VERSION 1.21.1
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN wget https://mirrors.aliyun.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz -O go.tgz \
    && tar -C /usr/local -xzf go.tgz \
    && rm go.tgz \
    && mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" \
    && go version
ENV GOPROXY https://goproxy.cn,direct
COPY . .
RUN bash build.sh

FROM debian:12
WORKDIR /app
COPY --from=builder /build/output .

EXPOSE 8888
CMD ["bash","bootstrap.sh"]
