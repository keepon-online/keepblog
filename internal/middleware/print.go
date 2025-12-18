package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 自定义一个结构体，实现 gin.ResponseWriter interface
type printResponseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w *printResponseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func PrintMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		//请求日期
		requestData := start.Format(time.DateTime)
		//请求接口路由
		path := c.Request.URL
		//请求方式
		requestMethod := c.Request.Method
		//请求 header
		//requestHeader := c.Request.Header
		//请求体 body
		requestBody := ""
		b, err := c.GetRawData()
		if err != nil {
			requestBody = "failed to get request body"
		} else {
			requestBody = string(b)
		}
		println(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
		//请求 host
		host := c.Request.Host
		//协议及版本
		schema := c.Request.Proto
		//客户端IP
		ip := c.ClientIP()

		writer := &printResponseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer

		c.Next()

		//请求耗时
		cost := time.Since(start).Milliseconds()
		//响应状态码
		responseStatus := c.Writer.Status()
		//响应 header
		//responseHeader := c.Writer.Header()
		//响应体大小
		responseBodySize := c.Writer.Size()
		//响应体 body
		//responseBody := writer.b.Bytes()
		info := PrintInfo{
			Time:             requestData,
			Path:             path.Path,
			Method:           requestMethod,
			Ip:               ip,
			Host:             host,
			Schema:           schema,
			Cost:             float32(cost) / 1000,
			Status:           responseStatus,
			RequestBody:      b,
			ResponseBodySize: responseBodySize,
			ResponseStatus:   responseStatus,
			ResponseBody:     writer.b.String(),
		}
		bbyt, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(string(bbyt))
	}
}

type PrintInfo struct {
	Time             string
	Path             string
	Method           string
	Ip               string
	Host             string
	Schema           string
	RequestHeader    http.Header
	Cost             float32
	Status           int
	RequestBody      []byte
	ResponseBodySize int
	ResponseHeader   http.Header
	ResponseStatus   int
	ResponseBody     string
}
