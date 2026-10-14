package link

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"syscall"
	"time"

	"gitee.com/jieepre/keepblog/internal/model"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
)

// safeDialer 建立连接时校验对端 IP：回环/私网/链路本地/未指定/组播地址一律拒绝，
// 防止友链检测接口被用作内网探测（SSRF）。在 Dial 层校验而非仅预解析域名，
// 可同时挡住 DNS 重绑定。
var safeDialer = &net.Dialer{
	Timeout: 5 * time.Second,
	Control: func(network, address string, _ syscall.RawConn) error {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return err
		}
		ip := net.ParseIP(host)
		if ip == nil {
			return fmt.Errorf("无法解析目标地址")
		}
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
			return fmt.Errorf("目标地址不在允许范围内")
		}
		return nil
	},
}

type Handler struct {
	*core.Context
}

func (h *Handler) SaveLink(c *gin.Context) {
	category := model.FriendLink{}
	err := c.ShouldBindJSON(&category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Save(category)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateLink(c *gin.Context) {
	link := model.FriendLink{}
	err := c.ShouldBindJSON(&link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Update(link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) UpdateLinkState(c *gin.Context) {
	link := model.FriendLink{}
	err := c.ShouldBindJSON(&link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	err = h.Service.LinkService.Update(link)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DeleteLink(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.ParseInt(param, 0, 64)

	err := h.Service.LinkService.Delete(int(id))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, nil)
}
func (h *Handler) DetailLink(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.ParseInt(param, 0, 64)
	linkInfo, err := h.Service.LinkService.GetLink(uint32(id))
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, linkInfo)
}
func (h *Handler) GetLinkList(c *gin.Context) {
	links, err := h.Service.LinkService.GetLinkList()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	result.Ok(c, links)
}

func (h *Handler) CheckLink(c *gin.Context) {
	targetUrl := c.Query("url")
	if targetUrl == "" {
		result.Error(c, "url 不能为空")
		return
	}
	u, err := url.Parse(targetUrl)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		result.Ok(c, gin.H{"status": 0, "msg": "无效链接"})
		return
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil // 不使用环境代理，避免绕过拨号校验
	transport.DialContext = safeDialer.DialContext
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
		// 不跟随跳转：跳转目标同样可能是内网地址
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
	if err != nil {
		result.Ok(c, gin.H{"status": 0, "msg": "无效链接"})
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		result.Ok(c, gin.H{"status": 0, "msg": "访问超时或无法连接"})
		return
	}
	defer resp.Body.Close()
	result.Ok(c, gin.H{"status": resp.StatusCode, "msg": resp.Status})
}

