package activity

import (
    "github.com/gin-gonic/gin"
    "time"
    "prosnav.com/wxserver/utils"
    "fmt"
    "crypto/sha1"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/wechat"
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/services/activityservice"
    "prosnav.com/wxserver/midwares"
    "prosnav.com/wxserver/modules/setting"
    "net/url"
)

type signForm struct {
    ActivityId int64 `form:"activityid" json"activityid,string"`
    CustName   string `form:"custname" json:"custname"`
    Mobile     string   `form:"mobile" json:"mobile"`
}

func SignIn(c *gin.Context) {
    var form signForm
    if err := c.Bind(&form); err != nil {
        c.JSON(500, nil)
        return
    }
    log.Debug("sinin form: %+v", form)
    sign := new(domain.SignIn)
    sign.ActivityId = form.ActivityId
    sign.Custname = form.CustName
    sign.Mobile = form.Mobile
    sess := midwares.GetSession(c)
    userid := sess.Get("user").(string)
    sign.UserId = userid
    ret, err := activityservice.SignIn(sign)
    if err != nil {
        panic(err)
        return
    }
    if ret == domain.SIGNED_IN {
        c.JSON(200, domain.SIGNED_IN)
        return
    }
    c.JSON(200, domain.OK)
}

type address struct {
    Latitude  float64   `json:"latitude" form:"latitude"`
    Longitude float64   `json:"longitude" form:"longitude"`
}

type result struct {
    Act      *domain.Activity    `json:"act"`
    UserName string      `json:"username"`
    Mobile   string      `json:"mobile"`
    Signed   bool    `json:"signed"`
}

func QueryActivity(c *gin.Context) {
    var addr address
    c.Bind(&addr)
    act, err := activityservice.QueryActivity(addr.Latitude, addr.Longitude)
    if err != nil {
        panic(err)
        return
    }
    ret := new(result)
    ret.Act = act
    if act == nil {
        c.JSON(200, ret)
        return
    }
    sess := midwares.GetSession(c)
    ret.UserName = sess.Get("username").(string)
    ret.Mobile = sess.Get("mobile").(string)
    userid := sess.Get("user").(string)
    signed, err := activityservice.Check(act.Id, userid)
    if err != nil {
        panic(err)
        return
    }
    ret.Signed = signed
    c.JSON(200, ret)
}

type actForm struct {
    Latitude     float64   `json:"latitude"`
    Longitude    float64   `json:"longitude"`
    Address      string `json:"address"`
    ActivityName string `json:"activityname"`
    Brief        string `json:"brief"`
    StartDate    int64  `json:"startdate,string"`
    EndDate      int64 `json:"enddate,string"`
    TagId        string `json:"tagid"`
    AgentId      int    `json:"agentid"`
}

func AddActivity(c *gin.Context) {
    var form actForm
    if err := c.Bind(&form); err != nil {
        panic(err)
    }
    act := new(domain.Activity)
    act.Latitude = form.Latitude
    act.Longitude = form.Longitude
    act.Address = form.Address
    act.ActivityName = form.ActivityName
    act.Brief = form.Brief
    act.StartDate = time.Unix(form.StartDate, 0)
    act.EndDate = time.Unix(form.EndDate, 0)
    act.TagId = form.TagId
    act.AgentId = form.AgentId
    if err := activityservice.AddActivity(act); err != nil {
        panic(err)
    }
    c.JSON(200, nil)
}

type initData struct {
    JsapiTicket string `json:"jsapi_ticket"`
    AppId       string `json:"appid"`
    Timestamp   int64  `json:"timestamp"`
    NonceStr    string `json:"nonceStr"`
    Signature   string `json:"signature"`
}

func InitData(c *gin.Context) {
    url := c.Query("url")
    log.Debug("location url: %s\n", url)
    initData := new(initData)
    initData.JsapiTicket = wechat.JsTicket
    initData.AppId = wechat.Appid
    initData.Timestamp = time.Now().Unix()
    initData.NonceStr = utils.RandomStr(32)
    toSignStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", initData.JsapiTicket, initData.NonceStr, initData.Timestamp, url)
    initData.Signature = fmt.Sprintf("%x", sha1.Sum([]byte(toSignStr)))
    c.JSON(200, initData)
}

func Authorize(c *gin.Context) {
    authUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=activity&scope=snsapi_base&#wechat_redirect",
        setting.Cfg.Section("wechat").Key("CORP_ID").String(),
        url.QueryEscape(setting.Cfg.Section("wechat").Key("REDIRECT_URI").String()))
    c.Redirect(301, authUrl)
}