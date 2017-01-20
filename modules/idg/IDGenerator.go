package idg

import (
    "prosnav.com/wxserver/modules/setting"
    "prosnav.com/wxserver/modules/log"
)

var (
    inst *IdWorker
)

func Id() (int64, error) {
    return inst.NextId()
}

func Ids(num int) ([]int64, error) {
    return inst.NextIds(num)
}

func init() {
    log.Debug("idg module start initializing.")
    datacenterId := setting.Cfg.Section("idg").Key("DATACENTER").MustInt64(2)
    worker, err := NewIdWorker(int64(1), datacenterId, twepoch)
    if err != nil {
        panic(err)
    }

    inst = worker
    log.Debug("idg module initialize successfully.")
}
