package oss

import (
	"bufio"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"gitee.com/jieepre/go-site/config"
	"gitee.com/jieepre/go-site/pkg"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/pkg/errors"
)

const (
	endpoint        string = "127.0.0.1:9000"
	accessKeyID     string = "root"
	secretAccessKey string = "123456abc"
	useSSL          bool   = false
)

var (
	client *minio.Client
	err    error
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

func createBucket(bucketName string) {
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-south-1", ObjectLocking: false})
	if err != nil {
		slog.Errorf("创建bucket错误: %s", err.Error())
		exists, _ := client.BucketExists(ctx, bucketName)
		if exists {
			slog.Errorf("bucket: %s已经存在", bucketName)
		}
	}
}

// FileUploader 上传
func FileUploader(bucketName, objectName, contextType string, size int64, data io.Reader) error {
	object, err := client.PutObject(context.Background(), bucketName, objectName, data, size, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		slog.Errorf("上传失败：%s", err.Error())
		return errors.New(err.Error())
	}
	slog.Infof("Successfully uploaded %s of size %d\n", objectName, object.Size)
	return nil
}

// GetFileUrl 获取文件url
func GetFileUrl(bucketName string, fileName string, expires time.Duration) string {
	//time.Second*24*60*60
	reqParams := make(url.Values)
	presignedURL, err := client.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		slog.Errorf(err.Error())
		return ""
	}
	return fmt.Sprintf("%s", presignedURL)
}

func lifecycleSet() {
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "expire-bucket",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: 365,
			},
		},
	}
	_ = client.SetBucketLifecycle(ctx, "mymusic", config)

}

func lifecycleGet() {
	bucketLifecycle, _ := client.GetBucketLifecycle(ctx, "mymusic")
	localLifecycleFile, _ := os.Create("lifecycle.json")
	defer localLifecycleFile.Close()

	encoder := xml.NewEncoder(localLifecycleFile)
	encoder.Indent("  ", "    ")
	encoder.Encode(bucketLifecycle)

}

func FileCopy() {
	src := minio.CopySrcOptions{
		Bucket:             "test2",
		Object:             "messages",
		MatchModifiedSince: time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC),
	}
	dst := minio.CopyDestOptions{
		Bucket: "mymusic",
		Object: "messages",
	}
	object, _ := client.CopyObject(ctx, dst, src)
	log.Printf("Copied %s, successfully to %s - UploadInfo %v\n", dst, src, object)
}

func FilesDelete() {
	bucketName := "mymusic"
	objectName := "audit.log"
	//删除一个文件
	_ = client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: true})

	//批量删除文件
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		options := minio.ListObjectsOptions{Prefix: "test", Recursive: true}
		for object := range client.ListObjects(ctx, bucketName, options) {
			if object.Err != nil {
				log.Println(object.Err)
			}
			objectsCh <- object
		}
	}()
	client.RemoveObjects(ctx, objectName, objectsCh, minio.RemoveObjectsOptions{})
}

func down() error {
	conf := config.Get().Minio

	imageUrl := pkg.GetPixabayImage()
	// 创建一个到远程图片的请求
	response, err := http.Get(imageUrl)
	if err != nil {
		slog.Errorf("Error while downloading: %s", err.Error())
		return errors.New(err.Error())
	}
	defer response.Body.Close()

	// 检查响应状态码是否为 200
	if response.StatusCode != http.StatusOK {
		slog.Errorf("Error while downloading %d", response.StatusCode)
		return errors.New(fmt.Sprintf("Error while dowload: HTTP status code:%d", response.StatusCode))
	}

	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(response.Body, 32*1024)
	file, err := os.Create("./cover.jpg")
	if err != nil {
		slog.Errorf(err.Error())
		return err
	}
	defer func() {
		file.Close()
	}()

	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	written, err := io.Copy(writer, reader)
	stat, err := os.Stat("./cover.jpg")
	if err != nil {
		slog.Errorf("获取文件大小失败：%s", err.Error())
		return errors.New("获取文件失败" + err.Error())
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return errors.New("创建新的UUID错误")
	}
	filename := newUUID.String() + path.Ext(stat.Name())
	buffer := make([]byte, 512)
	// Open File
	f, err := os.Open("./cover.jpg")
	if err != nil {
		slog.Errorf(err.Error())
		return err
	}
	defer f.Close()
	_, err = f.Read(buffer)
	if err != nil {
		slog.Errorf("获取文件类型失败：%s", err.Error())
		return errors.New("获取文件类型错误" + err.Error())
	}
	contentType := http.DetectContentType(buffer)
	err = FileUploader(conf.BucketName, filename, contentType, written, f)
	if err != nil {
		return err
	}
	return nil
}
