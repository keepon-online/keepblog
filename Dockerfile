FROM golang:1.19.9
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
COPY site.db /app/site.db
COPY ip2region.xdb /app/ip2region.xdb
COPY config/config.yaml /app/config/config.yaml
RUN go mod tidy
RUN go env && go build -o site .
FROM ubuntu:20.04
WORKDIR /app
COPY --from=0 /app/site .
COPY --from=0 /app/site.db .
COPY --from=0 /app/ip2region.xdb .
COPY --from=0 /app/config/config.yaml /app/config/config.yaml
EXPOSE 8000
EXPOSE 8589
EXPOSE 8890
ENTRYPOINT ./site
