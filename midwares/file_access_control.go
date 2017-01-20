package midwares
import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/utils"
)

func FileAccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        _, ok := c.Get("Security")
        if ok {
            c.Next()
            return
        }
        sess := GetSession(c)
        userid := sess.Get("user").(string)
        productCode := c.Params.ByName("tagCode")
        if isAccessable(userid, productCode) {
            c.Next()
            return
        }

        c.Redirect(304, "/wxqyh/v1/error/access_error.html")
    }
}

const productCount = `select count(*) count from t_wc_user where userid = $1 and array[$2]::varchar[] <@ tags`

type Count struct {
    Num int64 `xorm:"count"`
}


func isAccessable(userid, productCode string) bool {
    var c []*Count
    err := utils.Engine.Sql(productCount, userid, productCode).Find(&c)
    if err != nil {
        panic(err)
    }
    return c[0].Num > 0
}