# IP 数据库构建与更新闭环

项目使用 [ip2region](https://github.com/lionsoul2014/ip2region) 的 xdb 数据库进行 IPv4 归属查询。整条链路：

```
每周日 UTC 02:00 GitHub Actions（.github/workflows/ipdb.yml）:
  ip2region 固定提交的源数据 → xdb_maker gen
  → bench 一致性校验 → 项目 smoke test
  → 发布 Release ipdb-v{日期}-r{序号}（归档，不可变）
  → 重建 Release ipdb-latest（滚动指针）

运行时每日 03:00（pkg/area.UpdateIPDB）:
  下载 ipdb-latest 的 xdb + SHA256SUMS → 哈希校验
  → 未变化则跳过；变化则加载校验 → 原子替换 → 查询器热切换
  （任一步失败保留旧库；绝不回退官方源）

冷启动（pkg/area.EnsureIPDB，本地无 xdb 文件时）:
  自管 Release → 失败则回退官方 gitee/github 直链（Warn 日志，版本未锁定）
```

设计原则：**运行时日常更新只认自管 Release（版本锁定 + SHA-256 校验），不跟踪上游 master**；官方直链仅作为冷启动兜底，避免新部署在 GitHub 不可达时完全没有 IP 归属能力。

## CI 流程

`.github/workflows/ipdb.yml` 每周日 UTC 02:00 自动运行，也可手动触发（workflow push 仅在本文件变更时触发）。流程固定使用 ip2region 指定提交，执行：

1. 构建官方 `maker/golang` 下的 `xdb_maker`；
2. 从 `data/ipv4_source.txt` 生成 IPv4 xdb；
3. `xdb_maker bench` 对源数据和生成结果做一致性校验；
4. 项目 Go binding 加载 xdb 执行 smoke test（`scripts/verify-ipdb.go`）；
5. 生成 `SHA256SUMS` 和 `metadata.json`；
6. 上传 GitHub Actions Artifact（保留 90 天，审计用）；
7. **校验全部通过后**发布两个 Release（`GITHUB_TOKEN`，`contents: write`）：
   - `ipdb-v{YYYYMMDD}-r{run_number}`：归档，不可变，供回滚；
   - `ipdb-latest`：滚动指针，先删旧 Release 与 tag 再重建。

流水线不提交 `data/ip2region.xdb` 到仓库、不触发 Docker 发版。

## 运行时配置

```yaml
ipdb:
  # 完整 xdb 下载 URL；留空 = GitHub Release ipdb-latest
  # 国内服务器可指向自建镜像（镜像目录需同时提供 ip2region_v4.xdb 与 SHA256SUMS）
  updateUrl: ""
```

- `SHA256SUMS` 始终从 `updateUrl` 同级目录推导；
- 支持热重载（改配置即生效，下一次 03:00 更新使用新源）；
- `data/ip2region.xdb` 不入库（`.gitignore`），由冷启动下载或日常更新维护。

## 回滚

1. 在 GitHub Releases 找到目标 `ipdb-v*` 归档，复制资产直链
   （形如 `https://github.com/keepon-online/keepblog/releases/download/ipdb-vYYYYMMDD-rN/ip2region_v4.xdb`）；
2. 配置 `ipdb.updateUrl` 指向该直链；
3. 等待每日 03:00 自动更新（或重启服务）即完成回滚；
4. 如需恢复自动更新，把 `updateUrl` 置空即可。

## 本地验证

验证已有 xdb：

```bash
go run ./scripts/verify-ipdb.go data/ip2region.xdb
```

本地生成需要先获取官方仓库并构建 maker：

```bash
git clone https://github.com/lionsoul2014/ip2region.git /tmp/ip2region
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

如果将第三方数据源合并到源文件，必须先确认其是否允许离线复制、格式转换、合并、商业使用和重新分发。代码许可证不自动授予第三方 IP 数据的再分发权。每个产物应记录数据源、版本、生成器版本、生成时间和 SHA-256（见 Release 附带的 `metadata.json`）。

官方生成器文档：<https://github.com/lionsoul2014/ip2region/tree/master/maker/golang>
