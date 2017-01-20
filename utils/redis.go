package utils
import (
    "prosnav.com/wxserver/modules/setting"
    "time"
    "gopkg.in/redis.v2"
)
/*******************EXAMPLE**********************
* Reference: https://godoc.org/gopkg.in/redis.v2
*
************************************************/


var Pool *redis.Client


func init() {
    sec := setting.Cfg.Section("redis")
    server := sec.Key("REDIS_SERVER").MustString(":6379")

    Pool = redis.NewClient(&redis.Options{
        Network: "tcp",
        Addr: server,

        DialTimeout:  5 * time.Second,
        ReadTimeout:  time.Second,
        WriteTimeout: time.Second,

        PoolSize:    50,
        IdleTimeout:time.Minute * 10,
    })

}

