FROM golang:1.21.4 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
COPY data/ip2region.xdb /app/data/ip2region.xdb
COPY config.yaml /app/config.yaml
RUN go mod tidy
RUN go env && go build -ldflags="-s -w"  -o site .
FROM ubuntu:22.04
RUN groupadd -r appuser && \
    useradd -r -g appuser appuser && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app/data && \
    chmod 755 /app/data \
WORKDIR /app
VOLUME /app/data
USER appuser
COPY --from=builder /app/site /app/site
COPY --from=builder /app/data/ip2region.xdb /app/data/ip2region.xdb
COPY --from=builder /app/config.yaml /app/config.yaml
ENV TZ=Asia/Shanghai
EXPOSE 8000 8589 8890
ENTRYPOINT ["./site"]
