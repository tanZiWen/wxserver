package oauth

import (
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/midwares"
    "prosnav.com/wxserver/modules/wechat"
    "prosnav.com/wxserver/services/loginservice"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/setting"
    "fmt"
    "net/url"
    "prosnav.com/wxserver/modules/results"
    "strings"
    "github.com/chanxuehong/rand"
    "github.com/chanxuehong/sid"
    "net/http"
    "github.com/chanxuehong/wechat/mp/user/oauth2"
)

type loginForm struct {
    Userid   string `form:"userid" json:"userid"`
    Password string    `form:"passwd" json:"passwd"`
}

var (
    appMap map[string]string
    adminUrl = setting.Cfg.Section("publish").Key("MANAGEMENT_URL").String()
)

func Login(c *gin.Context) {
    var form loginForm
    c.Bind(&form)
    log.Debug("user form :%v\n", form)
    admin := loginservice.Login(form.Userid, form.Password)
    sess := midwares.GetSession(c)
    defer sess.SessionRelease(c.Writer)
    sess.Set("user", admin.UserId)
    sess.Set("username", admin.Name)

    log.Debug("success")
    c.JSON(200, nil)
}

func Authorize(c *gin.Context) {
    authUrl := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=publish&scope=snsapi_base&#wechat_redirect",
        setting.Cfg.Section("wechat").Key("CORP_ID").String(),
        url.QueryEscape(setting.Cfg.Section("wechat").Key("REDIRECT_URI").String()))
    c.Redirect(301, authUrl)
}

type initForm struct {
    Code  string  `form:"code" json:"code"`
    State string `form:"state" json:"state"`
}

func InitAccount(c *gin.Context) {
    var form initForm
    c.Bind(&form)
    var result map[string]interface{}
    if err := wechat.QueryUserWithCode(form.Code, &result); err != nil || result["UserId"] == nil {
        panic(new(results.UnAuthorizedError))
    }
    userid := result["UserId"].(string)
    wechat.QueryUserInfo(userid, &result)
    log.Debug("Wechat user info: %v\n", result)
    sess := midwares.GetSession(c)
    sess.Set("user", userid)
    sess.Set("username", result["name"])
    if result["mobile"] == nil {
        sess.Set("mobile", "")
    } else {
        sess.Set("mobile", result["mobile"])
    }
    log.Debug("session object: %v", sess)
    sess.SessionRelease(c.Writer)
    c.Redirect(301, appMap[form.State])
}

const CookieName  = "sessionid"

var (
    MP_APP_ID = setting.Cfg.Section("wechat.public").Key("APP_ID").String()
    MP_APP_SECRET = setting.Cfg.Section("wechat.public").Key("APP_SECRET").String()
    MP_REDIRECT_URL = setting.Cfg.Section("wechat.public").Key("REDIRECT_URI").String()
    MP_SCOPE = setting.Cfg.Section("wechat.public").Key("SCOPE").String()

    oauth2Config = oauth2.NewOAuth2Config(
        MP_APP_ID,
        MP_APP_SECRET,
        MP_REDIRECT_URL,
        MP_SCOPE,
    )
)

func Authorizion(c *gin.Context) {
    userAgent := c.Request.Header.Get("user-agent")

    if strings.Contains(userAgent, "MicroMessenger") {
        state := string(rand.NewHex())
        sid := sid.New()

        sess := midwares.GetSession(c)
        defer sess.SessionRelease(c.Writer)
        sess.Set(sid, state)

        cookie := http.Cookie{
            Name: CookieName,
            Value: sid,
            HttpOnly: true,
        }

        http.SetCookie(c.Writer, &cookie)

        authCodeUrl := oauth2Config.AuthCodeURL(state, nil)

        http.Redirect(c.Writer, c.Request, authCodeUrl, http.StatusFound)
    }else {
        c.JSON(200, "请在微信端尝试!")
    }
}

func OauthMP(c *gin.Context) {
    code := c.Query("code")
    state := c.Query("state")

    cookie, err := c.Request.Cookie(CookieName); if err != nil {
        log.Error(3, "get request cookie error: %v", err)
        c.AbortWithStatus(500)
        return
    }

    sess := midwares.GetSession(c)
    defer sess.SessionRelease(c.Writer)

    savedState := sess.Get(cookie.Value);

    log.Info("request:", savedState, state)

    if state != savedState {
        log.Error(3, "invalid state", "")
        c.AbortWithStatus(500)
        return
    }

    oauth2Client := oauth2.Client{
        Config: oauth2Config,
    }

    token, err := oauth2Client.Exchange(code)

    if err != nil {
        log.Error(3, "get token by code error:%v", err)
        c.AbortWithStatus(500)
        return
    }

    log.Info("%+v\n", token)
}

func init() {
    appMap = map[string]string{
        "publish": setting.Cfg.Section("publish").Key("NEWS_URL").String(),
        "activity": setting.Cfg.Section("activity").Key("SIGN_IN").String(),
    }
}