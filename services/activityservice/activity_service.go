package activityservice

import (
    "math"
    "strings"
    "prosnav.com/wxserver/domain"
    "prosnav.com/wxserver/modules/setting"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/idg"
    "prosnav.com/wxserver/utils"
    "fmt"
    "prosnav.com/wxserver/modules/wechat"
)

func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
    radius := float64(6371000) // 6378137
    rad := math.Pi/180.0

    lat1 = lat1 * rad
    lng1 = lng1 * rad
    lat2 = lat2 * rad
    lng2 = lng2 * rad

    theta := lng2 - lng1
    dist := math.Acos(math.Sin(lat1) * math.Sin(lat2) + math.Cos(lat1) * math.Cos(lat2) * math.Cos(theta))

    return dist * radius
}

const select_activity = `select id, longitude, latitude, address, activityname, brief, tagid, agentid from p_activity where now() between startdate and enddate + interval '1 day' and del = ?`

func filterByDistance(acts []domain.Activity, lat1, lng1 float64) *domain.Activity {
    var PRECISION = setting.Cfg.Section("distance").Key("PRECISION").MustFloat64(200.0)
    for _, act := range acts {
        dis := EarthDistance(lat1, lng1, act.Latitude, act.Longitude)
        log.Debug("Distance between (%f, %f) and (%f, %f) is %f\n", lat1, lng1, act.Latitude, act.Longitude, dis)
        if dis < PRECISION {
            return &act
        }
    }
    return nil
}

func QueryActivity(lat1, lng1 float64) (*domain.Activity, error) {
    var acts []domain.Activity
    if err := utils.Engine.Sql(select_activity, false).Find(&acts); err != nil {
        log.Error(4, "%v", err)
        return nil, err
    }
    log.Debug("Activities: %v\n", acts)
    return filterByDistance(acts, lat1, lng1), nil
}
func Check(actId int64, userid string)(bool, error) {
    var signIn domain.SignIn
    return utils.Engine.Where("activityid = ?", actId).And("userid = ?", userid).Get(&signIn)
}


func AddActivity(act *domain.Activity) error {
    if act.Id == 0 {
        act.Id, _ = idg.Id()
    }
    if _, err := utils.Engine.Insert(act); err != nil {
        log.Error(4, "%v", err)
        return err
    }
    return nil
}

func SignIn(signin *domain.SignIn) (string, error) {
    if signin.Id == 0 {
        signin.Id, _ = idg.Id()
    }
    sess := utils.Engine.NewSession()
    sess.Begin()
    if _, err := sess.Insert(signin); err != nil {
        log.Error(4, "%v", err)
        if strings.Contains(err.Error(), "unique") {
            return domain.SIGNED_IN, nil
        }
        return "", err
    }
    act := new(domain.Activity)
    if _, err := utils.Engine.Where("id = ?", signin.ActivityId).Get(act); err != nil {
        sess.Rollback()
        log.Debug("Query activity by id failed: %v\n", err)
        panic(err)
    }
    if err := wechat.AddMemberToTag(signin.UserId, act.TagId); err != nil {
        sess.Rollback()
        log.Debug("add tag failed: %v\n", err)
        panic(err)
    }
    msg := fmt.Sprintf(`{"touser": "%s","msgtype": "image","agentid": %d,"image": {"media_id": "%s"},"safe":"0"}`, signin.UserId, act.AgentId, act.MediaId)
    if err := wechat.SendMsg(msg); err != nil {
        sess.Rollback()
        log.Debug("Send message failed: %v\n", err)
        panic(err)
    }
    sess.Commit()
    return domain.OK, nil
}
