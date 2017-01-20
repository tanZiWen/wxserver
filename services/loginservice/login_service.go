package loginservice
import (
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/utils"
    "prosnav.com/wxserver/modules/results"
    "prosnav.com/wxserver/modules/log"
    "strings"
)

func Login(userid, passwd string) *domain.Admin {
    admin := new(domain.Admin)
    admin.UserId = userid
    has, err := utils.Engine.Get(admin)
    if err != nil {
        log.Error(3, "Query user data failed, detail: %v\n", err)
        panic(results.NewBusinessError("5002"))
    }
    if !has {
        panic(results.NewBusinessError("5002"))
    }
    salt := strings.Split(admin.Passwd, "_")[2]
    pswd, err := utils.EncryptPassword(passwd, salt)
    if err != nil {
        panic(err)
    }
    if pswd != admin.Passwd {
        panic(results.NewBusinessError("5003"))
    }

    return admin
}