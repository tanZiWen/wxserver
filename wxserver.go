package main



import (
    _"prosnav.com/wxserver/utils"
    "github.com/gin-gonic/gin"
    "prosnav.com/wxserver/modules/setting"
    "prosnav.com/wxserver/modules/mail"
    "prosnav.com/wxserver/controllers"
    "prosnav.com/wxserver/midwares"
    "prosnav.com/wxserver/modules/results"
        "fmt"
)


func ginDemo() {
    root := gin.New()
    appContext := fmt.Sprintf("/%s", setting.AppName)
    rootGroup := root.Group(appContext)
    version := fmt.Sprintf("/%s", setting.AppVer)

    verGroup := rootGroup.Group(version,
                            //log.Logger(),
                            midwares.SecurityAccess(),
                            midwares.AccessControl(),
                            gin.Recovery(),
                            results.ErrorHandler(),
                            //ginpongo2.Pongo2(),
    )

    controllers.RegisterHandlers(verGroup)
    root.Run(setting.ListeningPort)
}


func init() {
    mail.NewMailerContext()
}


func main() {
    ginDemo()
}




