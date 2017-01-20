package midwares

import (
	"github.com/astaxie/beego/session"
	"prosnav.com/wxserver/modules/setting"
	"github.com/gin-gonic/gin"
	_"github.com/astaxie/beego/session/redis"
	log "github.com/Sirupsen/logrus"
)

/**************************EXAMPLE**************************

    sess := GetSession(c)
    defer sess.SessionRelease(c.Writer)
    sess.Set("username", time.Now().Format(timeLayout))

************************************************************/

var globalSessions *session.Manager

func init() {
	var err error
	sessionConfig := &session.ManagerConfig{}
	section := setting.Cfg.Section("session")

	err = section.MapTo(sessionConfig)

	if err != nil {
		log.Error("failed to map ini to session config")
		panic(err)
	}

	globalSessions, err = session.NewManager(section.Key("PROVIDER").MustString("memory"), sessionConfig)
	if err != nil {
		log.Error("failed to init session middleware")
		panic(err)
	}

	//go globalSessions.GC()
}

func GetSession(c *gin.Context) session.Store {
    session, err := globalSessions.SessionStart(c.Writer, c.Request)
    if err != nil {
        panic(err)
    }
    return session
}