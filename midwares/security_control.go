package midwares
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/modules/setting"
    "prosnav.com/wxserver/modules/log"
)

var securityKey string

func SecurityAccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.Request.Header.Get("Security-Key")
        if key == securityKey {
            log.Debug("Security key: %s", key)
            c.Set("Security", key)
        }

        c.Next()
    }
}

func init() {
    securityKey = setting.Cfg.Section("auth").Key("KEYSECRET").String()
}