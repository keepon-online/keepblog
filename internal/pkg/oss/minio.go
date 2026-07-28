package oss

import (
	"context"
	"io"
	"net/url"
	"time"

	"gitee.com/jieepre/go-site/config"
	"github.com/gookit/slog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

var (
	client *minio.Client
	ctx    = context.Background()
)

func init() {
	m := config.Get().Minio
	minioClient, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
		Secure: m.UseSSL})
	if err != nil {
		slog.Errorf("minio连接错误:%s", err.Error())
	}
	client = minioClient
}

// FileUploader 上传文件到指定 bucket。
func FileUploader(bucketName, objectName, contextType string, size int64, data io.Reader) error {
	object, err := client.PutObject(ctx, bucketName, objectName, data, size, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		slog.Errorf("上传失败：%s", err.Error())
		return errors.New(err.Error())
	}
	slog.Infof("Successfully uploaded %s of size %d\n", objectName, object.Size)
	return nil
}

// GetFileUrl 获取文件预签名下载 URL。
func GetFileUrl(bucketName string, fileName string, expires time.Duration) string {
	reqParams := make(url.Values)
	presignedURL, err := client.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		slog.Errorf(err.Error())
		return ""
	}
	return presignedURL.String()
}
