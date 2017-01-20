package midwares
import (
    "github.com/gin-gonic/gin"
    "strings"
    "prosnav.com/wxserver/modules/log"
)

func AccessControl() gin.HandlerFunc {
    return func(c *gin.Context) {
        _, ok := c.Get("Security")
        if ok {
            c.Next()
            return
        }
        origin := c.Request.Header.Get("origin")
        log.Debug("origin: %s", origin)
        if origin == "" || strings.HasSuffix(origin, "prosnav.com") {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Next()
            return
        }

        c.AbortWithStatus(405)

    }
}
