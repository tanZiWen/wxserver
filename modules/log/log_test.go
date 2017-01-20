package log
import (
    "testing"
    "github.com/stretchr/testify/suite"
    "github.com/siddontang/go-log/log"
)

type staticStruct struct {
    suite.Suite
}

func (suite *staticStruct) SetupTest() {
    log.SetLevel(log.LevelInfo)
    log.Trace("Hello %s", "小黄")
}

func (suite *staticStruct)TestHello() {
    log.Trace("Hello %s", "小黄")
    log.Debug("Hello %s", "小黄")
    log.Info("Hello %s", "小黄")
    log.Error("Hello %s", "小黄")
}

func (suite *staticStruct)TearDownTest() {
}

func TestStatic(t *testing.T) {
    suite.Run(t, new(staticStruct))
}