package common

import (
	"fmt"
	"gitee.com/jieepre/keepblog/config"
	"gitee.com/jieepre/keepblog/internal/pkg/core"
	"gitee.com/jieepre/keepblog/internal/pkg/oss"
	"gitee.com/jieepre/keepblog/pkg/result"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path"
)

type Handler struct {
	*core.Context
}

func (h *Handler) UploadImage(c *gin.Context) {
	file, _ := c.FormFile("file")
	contentType := file.Header.Get("Content-Type")
	newUUID, err := uuid.NewUUID()
	if err != nil {
		result.Error(c, err.Error())
		return
	}

	filename := newUUID.String() + path.Ext(file.Filename)

	open, err := file.Open()
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	minio := config.Get().Minio
	bucketName := minio.BucketName
	err = oss.FileUploader(bucketName, filename, contentType, file.Size, open)
	if err != nil {
		result.Error(c, err.Error())
		return
	}
	baseUrl := minio.ServerUrl
	url := fmt.Sprintf("%s/%s/%s", baseUrl, bucketName, filename)
	result.Ok(c, url)
}
