package utils

import (
    "github.com/go-xorm/xorm"
    "prosnav.com/wxserver/modules/setting"
    _"github.com/lib/pq"
    "database/sql"
)

var Engine  *xorm.Engine

func init() {
    sec := setting.Cfg.Section("database")
    var (
        err error
    )
    Engine, err = xorm.NewEngine(sql.Drivers()[0], sec.Key("DSN").String())
    if err != nil {
        panic(err)
    }
    Engine.SetMaxOpenConns(sec.Key("MAX_CONNECTION").MustInt(50))
    Engine.SetMaxIdleConns(sec.Key("MAX_IDLE_CONNECTION").MustInt(50))
    Engine.ShowSQL()
}
