package log

import (
    "github.com/gogits/gogs/modules/log"
    "github.com/gin-gonic/gin"
    "time"
    "prosnav.com/wxserver/modules/setting"
)

func Trace(format string, v ...interface{}) {
    log.Trace(format, v...)
}

func Debug(format string, v ...interface{}) {
    log.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
    log.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
    log.Warn(format, v...)
}

func Error(skip int, format string, v ...interface{}) {
    log.Error(skip, format, v...)
}

func Critical(skip int, format string, v ...interface{}) {
    log.Critical(skip, format, v...)
}

func Fatal(skip int, format string, v ...interface{}) {
    log.Fatal(skip, format, v...)
}

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        end := time.Now()
        latency := end.Sub(start)
        clientIp := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        appName := setting.AppName
        log.Info("[%s] %d  %v | %s %s %s\n%s", appName, statusCode, latency,
                clientIp, method, c.Request.URL.Path, c.Errors.String())
    }
}
