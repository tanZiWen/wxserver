package wechat
import (
    "time"
    "fmt"
    "prosnav.com/wxserver/modules/setting"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/utils"
    "net/http"
    "strings"
    "encoding/json"
    "prosnav.com/wxserver/modules/results"
    "errors"
)


var (
    Appid, appSecret, AccessToken, JsTicket string
)

type WcResponse struct {
    AccessToken string	`json:"access_token"`
    ExpiresIn  int    `json:"expires_in"`
}

type (
    WechatToken struct {
        AccessToken  string `form:"access_token" json:"access_token"`
        ExpiresIn    int    `from:"expires_in" json:"expires_in"`
        Openid       string `from:"openid" json:"openid"`
        RefreshToken string `from:"refresh_token" json:"refresh_token"`
        ErrCode      int    `form:"errcode" json:"errcode"`
        ErrMsg       string `form:"errmsg" json:"errmsg"`
    }

    JsApiTicket struct {
        ErrCode int `json:"errcode"`
        ErrMsg  string `json:"errmsg"`
        Ticket  string `json:"ticket"`
        ExpiresIn int  `json:"expires_in"`
    }

    UserInfo struct {
        OpenId     string   `form:"openid" json:"openid"`
        NickName   string   `form:"nickname" json:"nickname"`
        Sex        int      `form:"sex" json:"sex"`
        Province   string   `form:"province" json:"province"`
        City       string   `form:"city" json:"city"`
        Country    string   `form:"country" json:"country"`
        HeadImgurl string   `form:"headimgurl" json:"headimgurl"`
        Privilege  []string `form:"privilege" json:"privilege"`
        UnionId    string   `form:"unionid" json:"unionid"`
        ErrCode    int      `form:"errcode" json:"errcode"`
        ErrMsg     string   `form:"errmsg" json:"errmsg"`
    }
)

func QueryUserWithCode(code string, result interface{}) error {
    userUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s&agentid=%d", AccessToken, code, 2)
    return utils.FetchUrlWithJson(userUrl, result)
}

func QueryUserInfo(userid string, result interface{}) {
    userInfoUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s", AccessToken, userid)
    utils.FetchUrlWithJson(userInfoUrl, result)
    return
}

func AddMemberToTag(userid, tagId string) error {
    addMemberUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/tag/addtagusers?access_token=%s", AccessToken)
    data := fmt.Sprintf(`{"tagid": "%s", "userlist":["%s"]}`, tagId, userid)
    log.Debug("parameter data : %s", data)
    reader := strings.NewReader(data)
    resp, err := http.Post(addMemberUrl, "application/json", reader)
    if err != nil {
        log.Debug("Add person to tag error: %v\n", err)
        return results.NewBusinessError("5010")
    }
    var result map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&result); err != nil {
        log.Debug("Add person to tag error: %v\n", err)
        return results.NewBusinessError("5010")
    }
    if result["errcode"].(float64) != 0 {
        log.Debug("Add person to tag result: %v\n", result)
        return results.NewBusinessError("5010")
    }
    log.Debug("Insert person into tag result: %v", result)
    return nil
}

func SendMsg(data string) error {
    sendMsgUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", AccessToken)
    //data := fmt.Sprintf(`{"touser": "%s","msgtype": "text","agentid": %d,"text": {"content": "%s"}, "safe":"0"}`, userid, agentId, content)
    reader := strings.NewReader(data)
    resp, err := http.Post(sendMsgUrl, "application/json", reader)
    if err != nil {
        log.Debug("Send text message to user :%s failed. %v\n", err)
        return results.NewBusinessError("5011")
    }
    var result map[string]interface{}
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&result); err != nil {
        log.Debug("Send text message to user :%s failed. %v\n", err)
        return results.NewBusinessError("5011")
    }
    if result["errcode"].(float64) != 0 {
        log.Debug("Send text message to user result :%v.\n", result)
        return results.NewBusinessError("5011")
    }
    log.Debug("Send message to user result: %v", result)
    return nil
}

func refreshToken() error {
    var token WechatToken
    accessTokenUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", Appid, appSecret)
    if err := utils.FetchUrlWithJson(accessTokenUrl, &token); err != nil {
        log.Error(4, "%v", err)
        return err
    }
    log.Debug("%+v", token)
    AccessToken = token.AccessToken
    return nil
}

func refreshJsTicket() error {
    jsApiUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=%s", AccessToken)
    var ticket JsApiTicket
    if err := utils.FetchUrlWithJson(jsApiUrl, &ticket); err != nil {
        log.Error(4, "%v", err)
        return err
    }
    if ticket.ErrCode != 0 {
        log.Error(4, "Request js api ticket failed, errmsg: %s", ticket.ErrMsg)
        return errors.New("Request js api ticket failed.")
    }
    JsTicket = ticket.Ticket
    return nil
}

func refresh() {
    for i := 0; i < 3; i ++ {
        if err := refreshToken(); err == nil {
            break
        }
        time.Sleep(3 * time.Second)
    }
    for i := 0; i < 3; i ++ {
        if err := refreshJsTicket(); err == nil {
            break
        }
        time.Sleep(3 * time.Second)
    }
}

func run() {
    refresh()
    tiker := time.NewTicker(7100 * time.Second)
    for {
        select {
        case <- tiker.C: refresh()
        }
    }
}

func init() {
    Appid = setting.Cfg.Section("wechat").Key("CORP_ID").String()
    appSecret = setting.Cfg.Section("wechat").Key("CORP_SECRET").String()
    go run()
}