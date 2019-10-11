package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"stockChina/common"
	"stockChina/conf"
	"stockChina/handler"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

const (
	TensentBaseUrl = "http://data.gtimg.cn/flashdata/hushen/daily/"
	SinaBaseUrl    = "http://hq.sinajs.cn/list="
)

// HandleDownloadFile 下载文件
func HandleDownloadFile(c *gin.Context) {
	fileList := []string{
		"images/1.jpeg",
		"images/2.jpeg",
		"images/letters.config",
	}
	//保留原来文件的结构
	err := ZipFiles("./test.zip", fileList)
	if err != nil {
		fmt.Println(err)
	}
	bytes, _ := ioutil.ReadFile("./test.zip")
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=hello.zip")
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(bytes)))
	c.Writer.Write(bytes)
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 优化压缩
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func setupRouter() *gin.Engine {

	r := gin.New()
	r.Use(common.Logger(), gin.Recovery())

	resolver := handler.TensentStockURLResolver{
		ExternalUrl: TensentBaseUrl,
	}
	HTTPClientReverseProxy := handler.NewHTTPClientReverseProxy(100*time.Second, 10, 10)
	GetKline := handler.MakeHandlerWrapper(handler.GetKline, resolver, HTTPClientReverseProxy)

	audit := r.Group("/v1")
	{
		audit.GET("/kline", GetKline)
	}

	r.GET("/_ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pretty double Qing !")
	})

	r.GET("download", HandleDownloadFile)

	return r
}

func main() {
	conf.InitEnvConfig()
	common.InitLogger()
	r := setupRouter()
	log.Info("Hello,I am stock sniper!Biu Biu Biu ~")
	r.Run(":8080")
}
