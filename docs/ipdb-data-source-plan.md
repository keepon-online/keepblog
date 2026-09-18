# IP 数据源扩充与评估规划

> 状态：规划中。当前生产主库仍为官方 ip2region IPv4 xdb，本规划不直接切换生产查询逻辑。

## 1. 背景与目标

当前项目使用 ip2region 提供 IP 归属查询，并已具备：

- 启动时缺失下载；
- 每日 03:00 运行时更新；
- 下载失败保留旧文件；
- xdb 加载校验和查询器热切换；
- 独立 CI 生成、校验和留存官方 xdb Artifact。

ip2region 官方没有公布可审计的全球覆盖率或城市级准确率。下一步目标不是立即替换它，而是建立可比较、可回滚的数据源评估流程，确认真实业务是否需要更高覆盖率或更细粒度定位。

## 2. 数据源候选

### 第一优先级：DB-IP Lite City

用于评估城市级结果，原因是它与当前 ip2region 的国家、省份、城市、ISP 类结果最接近。需要核实并记录：

- IPv4/IPv6 实际覆盖情况；
- CSV/MMDB 字段映射；
- CC BY 4.0 署名要求；
- 是否允许转换为 xdb、内部使用和再分发；
- 数据下载和更新频率。

### 第二优先级：IPinfo Lite

用于国家、洲、ASN 和组织归属评估。它支持 IPv4/IPv6 和每日更新，但不是城市级主库。转换或发布派生数据库时，需要遵守 CC BY-SA 4.0 的署名和 ShareAlike 要求。

### 辅助数据源：iptoasn / ip-location-db

用于补充 ASN、组织、云厂商和网络归属，不作为城市地理主库。接入前仍需记录来源、版本和许可证。

### 校验数据：RIR delegated stats

用于校验注册机构、国家代码、地址段状态和分配信息，不作为实际地理位置主库。

## 3. 分阶段实施

### 阶段一：建立真实样本基线

不修改生产查询逻辑，先统计脱敏后的真实查询样本：

- IPv4 和 IPv6 分开；
- 查询总数；
- 当前库命中数和未命中数；
- 国家、城市、ISP 结果为空的比例；
- 代理、CDN、云厂商等明显异常样本；
- 按时间和来源记录统计，不保存不必要的原始用户 IP。

样本统计需要限流或采样，避免日志膨胀和隐私风险。先确认主要缺口到底是 IPv6、海外地址、云出口，还是数据字段质量问题。

### 阶段二：候选库离线评估

新增独立评估任务，不接入 Web 请求路径：

```text
下载候选数据
    -> 校验许可证、版本和 SHA-256
    -> 转换为统一区间格式
    -> 检查非法、重复、重叠 IP 段
    -> 统计 IPv4/IPv6 网段和字段覆盖
    -> 对真实样本执行批量查询
    -> 与当前 ip2region 结果比较
    -> 上传评估报告
```

报告至少包含：

- 样本总数；
- 当前库命中率；
- 候选库命中率；
- 双方都命中的结果差异；
- 仅候选库命中的样本；
- 仅当前库命中的样本；
- IPv4/IPv6 分项结果；
- 国家、省份、城市、ISP、ASN 字段覆盖率；
- 数据文件大小、版本和 SHA-256；
- 评估时间和数据源版本。

评估报告只上传 CI Artifact，不自动修改 `master`、生产 xdb 或 Docker 镜像。

### 阶段三：双库灰度

只有评估结果和授权确认后，才增加候选库作为旁路查询：

```text
请求
  -> 当前 ip2region 主库返回结果
  -> 候选库旁路查询并记录差异
  -> 默认仍使用当前主库结果
```

旁路运行期间观察：

- 命中率；
- 结果差异率；
- 查询延迟；
- 内存占用；
- 数据更新失败率；
- 授权和部署成本。

### 阶段四：选择是否切换

满足以下条件后才考虑切换主库：

- 候选库在目标样本上的命中率有明确提升；
- 目标字段质量经过人工抽样确认；
- IPv4/IPv6 需求已明确；
- 数据源允许当前部署和发布方式；
- 有固定版本、SHA-256 和回滚产物；
- 性能和内存满足生产要求；
- 主库切换可通过配置或单次部署回滚。

如果提升不明显，继续使用 ip2region，仅保留现有更新和未命中统计。

## 4. 数据构建与发布要求

任何候选数据源都必须经过以下流水线：

```text
获取源数据
    -> 记录来源、版本、授权
    -> 格式转换
    -> IPv4/IPv6 分离
    -> 合法性、排序、重复和重叠检查
    -> 生成 xdb 或其他索引
    -> 官方/对应工具校验
    -> 固定 IP 样本 smoke test
    -> SHA-256
    -> 生成 metadata.json
    -> Artifact 或版本化 Release
```

`metadata.json` 至少记录：

```json
{
  "source": "数据源名称和 URL",
  "source_version": "数据源版本或抓取时间",
  "license": "许可证或合同标识",
  "generator": "生成器及版本",
  "generated_at": "UTC 时间",
  "sha256": "产物哈希",
  "size_bytes": 0
}
```

禁止：

- 未确认授权就合并第三方数据；
- 直接跟踪上游 `master` 作为生产版本；
- 下载失败时覆盖现有数据库；
- 把免费 API 批量抓取后当作可再分发数据库；
- 未经重叠和冲突处理直接拼接多个来源；
- 自动提交大体积二进制到业务主分支。

## 5. 当前推荐决策

短期保持：

```text
生产主库：官方 ip2region
独立评估：DB-IP Lite City
国家/ASN 对比：IPinfo Lite
ASN 辅助：iptoasn
注册信息校验：RIR delegated stats
```

如果只需要国家和 ASN，优先评估 IPinfo Lite；如果必须保持城市级字段，优先评估 DB-IP Lite City 或 GeoLite2 City。GeoLite2 的 EULA、更新和再分发义务需要单独审查，不能因为“免费”就默认允许公开发布转换后的 xdb。

## 6. 验收标准

本规划完成后的最小可交付结果：

- 有一条独立候选数据评估 workflow；
- 有一份可复现的候选数据版本和构建记录；
- 有覆盖率、命中率、字段覆盖率和差异率报告；
- 有授权与再分发结论；
- 生产主库未被自动替换；
- 评估失败不会影响主 CI、Docker 发版和线上查询；
- 保留当前 ip2region 版本作为可回滚主库。

## 7. 参考资料

- [ip2region 官方中文 README](https://github.com/lionsoul2014/ip2region/blob/master/README_zh.md)
- [KeepBlog 当前 IP 数据流水线](./ipdb.md)
- [IPinfo Lite](https://ipinfo.io/lite)
- [DB-IP Lite](https://db-ip.com/db/lite.php)
- [MaxMind GeoLite2](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data)
- [ip-location-db](https://github.com/sapics/ip-location-db)
- [RIR delegated stats](https://www.nro.net/about-the-nro/rir-statistics/)
