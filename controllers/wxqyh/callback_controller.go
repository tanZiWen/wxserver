package wxqyh
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/modules/wechat"
    "prosnav.com/wxserver/modules/log"
)


func Dispatch(c *gin.Context) {
    var msg wechat.TextRequestBody = c.Get("message").(wechat.TextRequestBody)

    if msg == nil {
        log.Debug("wechat message in:", msg)
    }

    c.JSON(201, nil)
}