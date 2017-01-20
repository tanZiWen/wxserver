package midwares
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/modules/results"
    "prosnav.com/wxserver/utils"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/setting"
    "strings"
)

var pass_urls []string

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        _, ok := c.Get("Security")
        if ok {
            c.Next()
            return
        }

        log.Debug("url: %s", c.Request.URL.Path)
        if utils.In(c.Request.URL.Path, pass_urls) {
            c.Next()
            return
        }
        user := authCheck(c)
        if user == nil {
            panic(new(results.UnAuthorizedError))
        }

        c.Next()
    }
}

func authCheck(c *gin.Context) interface{}{
    sess := GetSession(c)
    if v := sess.Get("user"); v != nil {
        return v
    }
    return nil
}

func init() {
    strs := setting.Cfg.Section("auth").Key("SKIP_URL").String()
    for _, s := range strings.Split(strs, ",") {
        pass_urls = append(pass_urls, strings.TrimSpace(s))
    }
}