package common

import (
	"io"
	"os"
	"time"

	"stockChina/conf"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)



// InitLogger init log options
func InitLogger() {
	var defaultWriter io.Writer
	toStdout := conf.EnvConf.Log.ToStdout
	logLevel := conf.EnvConf.Log.Level
	if toStdout {
		defaultWriter = os.Stdout
	} else {
		logSize := conf.EnvConf.Log.Size
		defaultWriter = &lumberjack.Logger{
			Filename:   "./stock.log",
			MaxSize:    logSize, // megabytes
			MaxBackups: 5,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		}
	}
	switch logLevel {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}

	logrus.SetOutput(defaultWriter)
}

//Logger return a self defined handler for log info
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logrus.Infof("|%3d |%13v |%15s |%s |%s |",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}





