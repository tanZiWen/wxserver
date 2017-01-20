package syncservice
import (
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/utils"
    "prosnav.com/wxserver/modules/log"
    "strings"
    "fmt"
)

const(
    user_exist = "select count(*) count from t_wc_user where userid = $1"
    user_update = "update t_wc_user set product = array[%s] where userid = $1"
    user_insert = "insert into t_wc_user(userid, product) values($1, array[%s])"
)

type Count struct {
    Num int64 `xorm:"count"`
}



func SyncUsers(users []*domain.User) map[string]bool{
    resultMap := make(map[string]bool, len(users))
    for _, user := range users {
        resultMap[user.UserId] = true
        var c []*Count
        err := utils.Engine.Sql(user_exist, user.UserId).Find(&c)
        if err != nil {
            resultMap[user.UserId] = false
            log.Error(3, "Query user %s failed\nError:%v", user.UserId, err)
            continue
        }

        params := make([]string, 0, len(user.Tags))
        for i := range user.Tags {
            params = append(params, fmt.Sprintf("$%v", i+2))
        }
        array := make([]interface{}, len(user.Tags) + 1)
        array[0] = user.UserId
        for i, v := range user.Tags {
            array[i + 1] = v
        }

        if c[0].Num > 0 {
            if _, err := utils.Engine.Exec(fmt.Sprintf(user_update, strings.Join(params, ", ")), array...); err != nil {
                log.Error(3, "Update user %s failed\nError:%v", user.UserId, err)
                resultMap[user.UserId] = false
            }
            continue
        }
        if _, err := utils.Engine.Exec(fmt.Sprintf(user_insert, strings.Join(params, ", ")), array...); err != nil {
            log.Error(3, "Insert user %s failed\nError:%v", user.UserId, err)
            resultMap[user.UserId] = false
        }

    }
    return resultMap
}

func SyncTags(tags []*domain.Tag) map[string]bool{
    resultMap := make(map[string]bool, len(tags))
    session := utils.Engine.NewSession()
    defer session.Close()
    session.Begin()
    for _, tag := range tags {
        resultMap[tag.TagCode] = true
        tmp := &domain.Tag{TagCode: tag.TagCode}
        exist, err := utils.Engine.Get(tmp); if err != nil {
            log.Error(3, "Get tag failed Error: %v", err)
            continue
        }
        if exist {
            if _, err := utils.Engine.Update(tag, tmp); err != nil {
                log.Error(3, "Update tag %s failed\nError:%v", tag.TagCode, err)
                resultMap[tag.TagCode] = false
            }
            continue
        }
        if _, err := utils.Engine.InsertOne(tag); err != nil {
            log.Error(3, "Insert product %s failed\nError:%v", tag.TagCode, err)
            resultMap[tag.TagCode] = false
        }
    }
    session.Commit()
    return resultMap
}
