package post

import "testing"

// TestFirstImageSrc og:image 兜底取正文首图。
func TestFirstImageSrc(t *testing.T) {
	cases := []struct {
		name, html, want string
	}{
		{"常规", `<p><img src="https://a.dev/1.png" loading="lazy"></p>`, "https://a.dev/1.png"},
		{"多图取第一张", `<img src="https://a.dev/x.jpg"><p>文字</p><img src="https://a.dev/y.jpg">`, "https://a.dev/x.jpg"},
		{"无图", `<p>纯文字正文</p>`, ""},
		{"空串", ``, ""},
	}
	for _, c := range cases {
		if got := firstImageSrc(c.html); got != c.want {
			t.Errorf("%s: firstImageSrc = %q, want %q", c.name, got, c.want)
		}
	}
}
