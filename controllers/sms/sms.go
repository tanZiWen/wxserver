package sms

import (
    "github.com/gin-gonic/gin"
    "prosnav.com/common/utils"
    "prosnav.com/wxserver/midwares"
    "prosnav.com/wxserver/modules/redis"
    "fmt"
    "time"
    "strings"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/yunpian"
)

const (
    smsKeyPrefix = "ps:sms:"
    mobile_or_code_empty  = "mobile_or_code_empty"
    code_timeout = "code_timeout"
    invalid_code = "invalid_code"
    sms_error = "sms_error"
)

func RefreshCode(c *gin.Context) {
    mobile := c.Query("mobile")
    code := utils.RandomSpecStr(3, utils.INT_CHARSET)
    sess := midwares.GetSession(c)
    sid := sess.SessionID()
    smsKey := fmt.Sprintf("%s%s:%s", smsKeyPrefix, sid, mobile)
    redis.RedisClient.SetEx(smsKey, 60 * time.Second, code)
    if err := yunpian.SendMsg(code, mobile); err != nil {
        log.Debug(err.Error())
        c.JSON(500, sms_error)
        return
    }
    sess.SessionRelease(c.Writer)
    c.JSON(200, nil)
}

type smsForm struct  {
    Mobile string `form:"mobile" json:"mobile"`
    Code   string `form:"code" json:"code"`
}



func ValidateCode(c *gin.Context) {
    var form smsForm
    c.Bind(&form)
    if strings.TrimSpace(form.Mobile) == "" || strings.TrimSpace(form.Code) == "" {
        c.JSON(500, mobile_or_code_empty)
        return
    }
    sess := midwares.GetSession(c)
    sid := sess.SessionID()
    smsKey := fmt.Sprintf("%s%s:%s", smsKeyPrefix, sid, form.Mobile)
    expectCode := redis.RedisClient.Get(smsKey).Val()
    log.Debug("Cached code: %s  form code: %s\n", expectCode, form.Code)
    if expectCode == ""{
        c.JSON(500, code_timeout)
        return
    }
    if expectCode != form.Code {
        c.JSON(500, invalid_code)
        return
    }
    c.JSON(200, nil)
}
