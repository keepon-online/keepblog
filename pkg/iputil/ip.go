package iputil

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// GetRealAddressByIP 返回 IP 的可读地址。
func GetRealAddressByIP(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "未知地址"
	}
	if isLocalIP(parsed) {
		return "服务器登录"
	}
	if isLANIP(parsed) {
		return "局域网"
	}
	return getLocation(ip)
}

func isLocalIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast()
}

func isLANIP(ip net.IP) bool {
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return ip4[0] == 10 ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
		(ip4[0] == 192 && ip4[1] == 168)
}

type locationResponse struct {
	Code int    `json:"code"`
	Addr string `json:"addr"`
}

func getLocation(ip string) string {
	url := "https://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	decoded, err := simplifiedchinese.GBK.NewDecoder().String(string(body))
	if err != nil {
		return ""
	}
	var result locationResponse
	if err := json.Unmarshal([]byte(decoded), &result); err != nil {
		return ""
	}
	if result.Code == 0 && result.Addr != "" {
		return result.Addr
	}
	return "未知地址"
}

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok || ipAddr.IP.IsLoopback() || !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return "", nil
}
