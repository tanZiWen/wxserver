package redis

import (
    "gopkg.in/redis.v2"
    "time"
    "encoding/gob"
    "bytes"
    "prosnav.com/wxserver/modules/log"
    "prosnav.com/wxserver/modules/setting"
)

/*******************EXAMPLE**********************
* Reference: https://godoc.org/gopkg.in/redis.v2
*
************************************************/

var (
    RedisClient *redis.Client
)

func init() {
    log.Debug("redis module start initializing.")
    server := setting.Cfg.Section("redis").Key("REDIS_SERVER").MustString(":6379")

    RedisClient = redis.NewClient(&redis.Options{
        Network: "tcp",
        Addr:    server,

        DialTimeout:  5 * time.Second,
        ReadTimeout:  time.Second,
        WriteTimeout: time.Second,

        PoolSize:    setting.Cfg.Section("redis").Key("POOL_SIZE").MustInt(50),
        IdleTimeout: time.Second * time.Duration(setting.Cfg.Section("redis").Key("IDEL_TIMEOUT").MustInt64(600)),
    })

    log.Debug("redis module initialize successfully.")
}

func SetExObj(key string, obj interface{}, t time.Duration) error {
    data, err := EncodeGob(obj)
    if err != nil {
        log.Debug("Encode gob failed: %v", err)
        return err
    }
    RedisClient.SetEx(key, t, string(data))
    return nil
}

func SetObj(key string, obj interface{}) error {
    data, err := EncodeGob(obj)
    if err != nil {
        log.Debug("Encode gob failed: %v", err)
        return err
    }
    RedisClient.Set(key, string(data))
    return nil
}

func GetObj(key string, obj interface{}) error {
    data := RedisClient.Get(key).Val()
    return DecodeGob([]byte(data), obj)
}

func EncodeGob(obj interface{}) ([]byte, error) {
    gob.Register(obj)
    buf := bytes.NewBuffer(nil)
    enc := gob.NewEncoder(buf)
    err := enc.Encode(obj)
    if err != nil {
        return []byte(""), err
    }
    return buf.Bytes(), nil
}

func DecodeGob(encoded []byte, obj interface{}) error {
    buf := bytes.NewBuffer(encoded)
    dec := gob.NewDecoder(buf)
    err := dec.Decode(&obj)
    if err != nil {
        return err
    }
    return nil
}
