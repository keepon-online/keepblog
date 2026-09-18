package area

// 查询样本基线统计（docs/ipdb-data-source-plan.md 阶段一）。
//
// 设计约束：
// - 纯内存原子计数，不保存任何原始用户 IP（隐私）；
// - Area() 是全站 IP 归属查询的唯一入口，统计在此收集；
// - 每日由定时任务输出 JSON 快照到日志并清零，日志即历史记录；
// - xdb 区域格式为"国家|区域|省份|城市|ISP"，字段占位值为空串或 "0"。

import (
	"encoding/json"
	"strings"
	"sync/atomic"

	"github.com/gookit/slog"
)

type baselineStats struct {
	Queries    atomic.Int64 // 总查询数
	NoDB       atomic.Int64 // 查询器未初始化（无库可用）
	Hits       atomic.Int64 // 返回了非空归属
	Misses     atomic.Int64 // 返回空/查询出错
	Private    atomic.Int64 // 内网/回环等保留地址（归属结果无意义）
	CloudOrCDN atomic.Int64 // 归属中含云厂商/CDN 关键字（代理出口常见）
	// 字段空置计数（分母为 Hits）
	CountryEmpty  atomic.Int64
	ProvinceEmpty atomic.Int64
	CityEmpty     atomic.Int64
	ISPEmpty      atomic.Int64
}

var stats baselineStats

// cloudCDNKeywords 归属串中出现即视为"云/CDN 出口"特征。
// 保守清单：宁缺毋滥，评估报告中只作异常样本提示，不作结论。
var cloudCDNKeywords = []string{
	"阿里云", "腾讯云", "华为云", "百度云", "京东云", "青云", "UCloud", "火山引擎",
	"亚马逊", "AWS", "Amazon", "微软", "Azure", "Microsoft",
	"谷歌", "Google", "Cloudflare", "Akamai", "Fastly",
	"DigitalOcean", "Linode", "Vultr", "OVH", "Hetzner",
	"机房", "数据中心", "IDC", "CDN",
}

// recordQuery 由 Area() 调用，region 为查询结果（可能为空）。
func recordQuery(intIP uint32, region string) {
	stats.Queries.Add(1)

	if isPrivateIPv4(intIP) {
		stats.Private.Add(1)
		return // 内网地址不参与字段统计，避免污染基线
	}
	if region == "" {
		stats.Misses.Add(1)
		return
	}
	stats.Hits.Add(1)

	if matchCloudCDN(region) {
		stats.CloudOrCDN.Add(1)
	}

	// 字段空置：国家|区域|省份|城市|ISP
	parts := strings.Split(region, "|")
	empty := func(i int) bool {
		return i >= len(parts) || parts[i] == "" || parts[i] == "0"
	}
	if empty(0) {
		stats.CountryEmpty.Add(1)
	}
	if empty(2) {
		stats.ProvinceEmpty.Add(1)
	}
	if empty(3) {
		stats.CityEmpty.Add(1)
	}
	if empty(4) {
		stats.ISPEmpty.Add(1)
	}
}

func recordNoDB() {
	stats.Queries.Add(1)
	stats.NoDB.Add(1)
}

func matchCloudCDN(region string) bool {
	for _, kw := range cloudCDNKeywords {
		if strings.Contains(region, kw) {
			return true
		}
	}
	return false
}

// isPrivateIPv4 判断保留/内网地址段。这些地址在 xdb 中是占位归属，
// 对评估真实数据覆盖率没有意义。
func isPrivateIPv4(ip uint32) bool {
	b1 := ip >> 24
	b2 := (ip >> 16) & 0xff
	switch {
	case b1 == 10, b1 == 127: // 内网A / 回环
		return true
	case b1 == 172 && b2 >= 16 && b2 <= 31: // 内网B
		return true
	case b1 == 192 && b2 == 168: // 内网C
		return true
	case b1 == 169 && b2 == 254: // 链路本地
		return true
	case b1 == 100 && b2 >= 64 && b2 <= 127: // CGNAT
		return true
	case b1 == 0, b1 == 255: // 本网络 / 广播
		return true
	}
	return false
}

// StatsSnapshot 返回当前统计快照（字段空置率以命中数为分母）。
func StatsSnapshot() map[string]any {
	queries := stats.Queries.Load()
	hits := stats.Hits.Load()
	rate := func(n int64) float64 {
		if hits == 0 {
			return 0
		}
		return float64(n) / float64(hits)
	}
	return map[string]any{
		"queries":        queries,
		"hits":           hits,
		"misses":         stats.Misses.Load(),
		"no_db":          stats.NoDB.Load(),
		"private":        stats.Private.Load(),
		"cloud_or_cdn":   stats.CloudOrCDN.Load(),
		"country_empty":  rate(stats.CountryEmpty.Load()),
		"province_empty": rate(stats.ProvinceEmpty.Load()),
		"city_empty":     rate(stats.CityEmpty.Load()),
		"isp_empty":      rate(stats.ISPEmpty.Load()),
	}
}

// ResetStats 清零（仅测试与快照输出后使用）。
func ResetStats() {
	stats = baselineStats{}
}

// LogStatsSnapshot 输出 JSON 快照到日志并清零，由每日定时任务调用。
// 日志即阶段一基线的历史记录，grep "ipdb_stats" 可回溯。
func LogStatsSnapshot() {
	snapshot := StatsSnapshot()
	data, err := json.Marshal(snapshot)
	if err != nil {
		slog.Warnf("ipdb 统计快照序列化失败: %v", err)
		return
	}
	slog.Infof("ipdb_stats %s", string(data))
	ResetStats()
}
