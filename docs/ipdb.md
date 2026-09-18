# IP 数据库构建流水线

项目使用 [ip2region](https://github.com/lionsoul2014/ip2region) 的 xdb 数据库进行 IPv4 归属查询。应用运行时会在每天 03:00 尝试更新；本流水线则独立负责从官方源数据重新生成、校验并留存 xdb 产物。

## CI 流程

`.github/workflows/ipdb.yml` 每周日 UTC 02:00 自动运行，也可以手动触发。流程固定使用 ip2region `v3.9.0` 的提交，执行：

1. 构建官方 `maker/golang` 下的 `xdb_maker`；
2. 从 `data/ipv4_source.txt` 生成 IPv4 xdb；
3. 使用 `xdb_maker bench` 对源数据和生成结果进行一致性校验；
4. 使用项目的 Go binding 加载 xdb 并执行 smoke test；
5. 生成 `SHA256SUMS` 和 `metadata.json`；
6. 上传 GitHub Actions Artifact，保留 90 天。

流水线不会自动提交 `data/ip2region.xdb`，也不会自动创建 Release，避免 IP 数据更新污染业务提交历史或触发 Docker 发版。审核通过后，可从 Artifact 下载并部署到 `/app/data/ip2region.xdb`，部署前应保留旧版本以便回滚。

## 本地验证

验证已有 xdb：

```bash
go run ./scripts/verify-ipdb.go data/ip2region.xdb
```

本地生成需要先获取官方仓库并构建 maker：

```bash
git clone --branch v3.9.0 https://github.com/lionsoul2014/ip2region.git /tmp/ip2region
cd /tmp/ip2region/maker/golang
go build -o xdb_maker .
./xdb_maker gen \
  --src=../../data/ipv4_source.txt \
  --dst=/tmp/ip2region_v4.xdb \
  --version=ipv4
./xdb_maker bench \
  --db=/tmp/ip2region_v4.xdb \
  --src=../../data/ipv4_source.txt \
  --version=ipv4
cd -
go run ./scripts/verify-ipdb.go /tmp/ip2region_v4.xdb
sha256sum /tmp/ip2region_v4.xdb
```

IPv6 数据需要使用 `ipv6_source.txt`、`--version=ipv6`，并由应用额外接入 IPv6 查询逻辑；当前项目流水线暂只构建 IPv4。

## 数据边界与授权

官方 README 没有公布全球覆盖率、国家覆盖率或城市级准确率，因此不能将该数据解释为所有 IP 的精确物理位置。VPN、代理、NAT、CDN、云出口和移动网络等地址可能存在偏差。

如果将第三方数据源合并到源文件，必须先确认其是否允许离线复制、格式转换、合并、商业使用和重新分发。代码许可证不自动授予第三方 IP 数据的再分发权。每个产物应记录数据源、版本、生成器版本、生成时间和 SHA-256。

官方生成器文档：<https://github.com/lionsoul2014/ip2region/tree/master/maker/golang>
