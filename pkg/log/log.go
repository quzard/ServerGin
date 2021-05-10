package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

/*
	案例
   //Info级别的日志
   Logger().WithFields(logrus.Fields{
       "name": "hanyun",
   }).Info("记录一下日志", "Info")
   //Error级别的日志
   log.Logger().WithFields(logrus.Fields{
       "name": "hanyun",
   }).Error("记录一下日志", "Error")
   //Warn级别的日志
   Logger().WithFields(logrus.Fields{
       "name": "hanyun",
   }).Warn("记录一下日志", "Warn")
   //Debug级别的日志
   Logger().WithFields(logrus.Fields{
       "name": "hanyun",
   }).Debug("记录一下日志", "Debug")

*/
var Log *logrus.Logger
func Logger()  {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/data/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		return
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			return
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return
	}

	//实例化
	Log = logrus.New()

	//设置输出
	Log.Out = src

	//设置日志级别
	Log.SetLevel(logrus.DebugLevel)

	//设置日志格式
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func LoggerToFile() (gin.HandlerFunc, error) {
	Logger()
	return func(c *gin.Context) {
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

		//日志格式
		Log.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}, nil
}
