package yunpian

import (
    "github.com/cJrong/YunPianSMS"
    "fmt"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/setting"
)

var apikey, company string

func SendMsg(code, mobile string) error {
    msg := fmt.Sprintf("【%s】您的验证码是%s", company, code)
    result, err := YunPianSMS.YunPianSMSSend(apikey, msg, mobile)
    log.Debug("Yunpian sms send result: %+v\n", result)
    return err
}

func init() {
    apikey = setting.Cfg.Section("sms").Key("YUNPIAN_APIKEY").String()
    company = setting.Cfg.Section("sms").Key("COMPANY").String()
}