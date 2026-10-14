package common

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/pkg/oss"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	*core.Context
}

// 上传图片限制：扩展名白名单 + 文件头真实类型校验 + 大小上限。
// 该接口曾因 JWT 白名单前缀而被匿名放行（可上传 HTML/SVG 同域托管），
// 现已回到管理鉴权组，此处再做内容层面的收敛。
const maxUploadImageSize = 10 << 20 // 10MB

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true, ".bmp": true,
}

func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		result.Error(c, "请选择要上传的图片文件")
		return
	}
	if file.Size > maxUploadImageSize {
		result.Error(c, "图片大小不能超过 10MB")
		return
	}
	ext := strings.ToLower(path.Ext(file.Filename))
	if !allowedImageExts[ext] {
		result.Error(c, "仅支持 jpg/jpeg/png/gif/webp/bmp 格式的图片")
		return
	}

	open, err := file.Open()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	defer func() { _ = open.Close() }()

	// 客户端自报的 Content-Type 不可信：以文件头识别真实类型，拒绝伪装成图片的其它内容
	head := make([]byte, 512)
	n, _ := io.ReadFull(open, head)
	contentType := http.DetectContentType(head[:n])
	if !strings.HasPrefix(contentType, "image/") {
		result.Error(c, "文件内容不是有效的图片")
		return
	}
	if _, err := open.Seek(0, io.SeekStart); err != nil {
		result.Error(c, err.Error())
		return
	}

	newUUID, err := uuid.NewUUID()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	filename := newUUID.String() + ext

	minio := config.Get().Minio
	bucketName := minio.BucketName
	if err := oss.FileUploader(bucketName, filename, contentType, file.Size, open); err != nil {
		result.Error(c, err.Error())
		return
	}
	baseUrl := minio.ServerUrl
	url := fmt.Sprintf("%s/%s/%s", baseUrl, bucketName, filename)
	result.Ok(c, url)
}
