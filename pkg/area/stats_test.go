package area

import (
	"encoding/binary"
	"net"
	"testing"
)

func ipToUint(t *testing.T, ip string) uint32 {
	t.Helper()
	parsed := net.ParseIP(ip).To4()
	if parsed == nil {
		t.Fatalf("非法测试 IP: %s", ip)
	}
	return binary.BigEndian.Uint32(parsed)
}

func TestRecordQueryClassification(t *testing.T) {
	ResetStats()
	t.Cleanup(ResetStats)

	// 完整字段命中
	recordQuery(ipToUint(t, "1.2.3.4"), "中国|华南|广东省|深圳市|电信")
	// 占位字段命中(国家有,省/市/ISP 为 0)
	recordQuery(ipToUint(t, "5.6.7.8"), "美国|0|0|0|0")
	// 云厂商出口
	recordQuery(ipToUint(t, "9.10.11.12"), "中国|0|北京市|北京市|阿里云")
	// 未命中
	recordQuery(ipToUint(t, "13.14.15.16"), "")
	// 内网地址
	recordQuery(ipToUint(t, "192.168.1.10"), "0|0|0|0|内网")

	s := StatsSnapshot()
	if s["queries"] != int64(5) {
		t.Errorf("queries = %v, want 5", s["queries"])
	}
	if s["hits"] != int64(3) {
		t.Errorf("hits = %v, want 3", s["hits"])
	}
	if s["misses"] != int64(1) {
		t.Errorf("misses = %v, want 1", s["misses"])
	}
	if s["private"] != int64(1) {
		t.Errorf("private = %v, want 1", s["private"])
	}
	if s["cloud_or_cdn"] != int64(1) {
		t.Errorf("cloud_or_cdn = %v, want 1", s["cloud_or_cdn"])
	}
	// 国家空置:3 次命中中 0 次
	if s["country_empty"] != 0.0 {
		t.Errorf("country_empty = %v, want 0", s["country_empty"])
	}
	// 省份空置:美国那次 + 阿里云那次(北京市有值? "北京市"非 0,只有美国 0)
	if s["province_empty"] != 1.0/3 {
		t.Errorf("province_empty = %v, want %v", s["province_empty"], 1.0/3)
	}
	// ISP 空置:美国一次;阿里云命中 ISP 有值
	if s["isp_empty"] != 1.0/3 {
		t.Errorf("isp_empty = %v, want %v", s["isp_empty"], 1.0/3)
	}
}

func TestRecordNoDB(t *testing.T) {
	ResetStats()
	t.Cleanup(ResetStats)

	recordNoDB()
	recordNoDB()

	s := StatsSnapshot()
	if s["queries"] != int64(2) || s["no_db"] != int64(2) {
		t.Errorf("queries=%v no_db=%v, want 2/2", s["queries"], s["no_db"])
	}
}

func TestIsPrivateIPv4(t *testing.T) {
	cases := map[string]bool{
		"10.0.0.1": true, "10.255.0.1": true,
		"172.16.0.1": true, "172.31.255.1": true, "172.32.0.1": false,
		"192.168.1.1": true, "192.169.1.1": false,
		"127.0.0.1":   true,
		"169.254.1.1": true,
		"100.64.0.1":  true, "100.63.0.1": false, "100.128.0.1": false,
		"0.1.2.3": true, "255.255.255.255": true,
		"8.8.8.8": false, "1.2.3.4": false,
	}
	for ip, want := range cases {
		if got := isPrivateIPv4(ipToUint(t, ip)); got != want {
			t.Errorf("isPrivateIPv4(%s) = %v, want %v", ip, got, want)
		}
	}
}

func TestSnapshotRatesWithoutHits(t *testing.T) {
	ResetStats()
	t.Cleanup(ResetStats)

	// 只有未命中时,空置率应为 0 而非除零
	recordQuery(ipToUint(t, "1.1.1.1"), "")
	s := StatsSnapshot()
	if s["city_empty"] != 0.0 {
		t.Errorf("无命中时 city_empty = %v, want 0", s["city_empty"])
	}
}

func TestLogStatsSnapshotResets(t *testing.T) {
	ResetStats()
	t.Cleanup(ResetStats)

	recordQuery(ipToUint(t, "8.8.4.4"), "美国|0|0|0|Google")
	LogStatsSnapshot()

	s := StatsSnapshot()
	if s["queries"] != int64(0) {
		t.Errorf("快照输出后未清零, queries = %v", s["queries"])
	}
}
