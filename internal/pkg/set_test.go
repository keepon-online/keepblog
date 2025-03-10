package pkg

import (
	"testing"
)

func TestName(t *testing.T) {
	api := "http://data.zz.baidu.com/urls?site=https://www.keepon.online&token=VzXKVoNrqnSAWINA"
	PushSite([]string{"https://www.keepon.online/post/KDa7MAqM1NPq1", "https://www.keepon.online/post/e6o4MzdpAaqX"}, api)
}
