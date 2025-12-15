FROM golang:1.24.11 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
COPY data/ip2region.xdb /app/data/ip2region.xdb
COPY config-example.yaml /app/config.yaml
RUN go mod tidy
RUN go env && go build -ldflags="-s -w"  -o site .
FROM ubuntu:22.04
RUN groupadd -r appuser && \
    useradd -r -g appuser appuser && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app/data

WORKDIR /app

COPY --from=builder /app/site /app/site
COPY --from=builder /app/config.yaml /app/default/config.yaml
COPY --from=builder /app/data/ip2region.xdb /app/default/ip2region.xdb

# 添加入口点脚本
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENV TZ=Asia/Shanghai
EXPOSE  8589

ENTRYPOINT ["/app/entrypoint.sh"]
