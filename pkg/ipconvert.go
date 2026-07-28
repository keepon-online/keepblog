package pkg

import (
	"encoding/binary"
	"net"
	"strings"
)

func Ip2long(ipstr string) uint32 {
	if ipstr == "::1" {
		ipstr = "127.0.0.1"
	}

	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	if isIpv4(ipstr) {
		ip = ip.To4()
		return binary.BigEndian.Uint32(ip)
	}
	return 0
}

// 判断是否ipv4
func isIpv4(ipStr string) bool {

	ip := net.ParseIP(ipStr)

	return ip != nil && strings.Contains(ipStr, ".")
}

func Long2ip(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
