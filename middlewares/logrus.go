package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func LoggerFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		file_dir := "logs/" + "gin_project.log"

		//写入文件
		src, err := os.OpenFile(file_dir, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			fmt.Println("err", err)
		}

		//实例化
		logger := logrus.New()

		//设置输出
		logger.Out = src

		//设置日志级别
		logger.SetLevel(logrus.DebugLevel)

		//设置日志格式，格式化时间
		logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)

	}

}
