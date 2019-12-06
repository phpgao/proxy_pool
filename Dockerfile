FROM golang:1.13.4-alpine3.10 AS build
COPY . /app
WORKDIR /app
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
#    apk add git && \
#    git clone https://github.com/phpgao/proxy_pool && \
RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy_pool


FROM alpine:3.10
COPY --from=build /app/proxy_pool /app/proxy_pool
WORKDIR /app
ENTRYPOINT ["/app/proxy_pool"]