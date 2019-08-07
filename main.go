package main

import (
	"net/http"
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

	return r
}

func main() {
	conf.InitEnvConfig()
	common.InitLogger()
	r := setupRouter()
	log.Info("Hello,I am stock sniper!Biu Biu Biu ~")
	r.Run(":8080")
}
