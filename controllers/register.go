package controllers

import (
	"github.com/gin-gonic/gin"
	"prosnav.com/wxserver/controllers/oauth"
	"prosnav.com/wxserver/controllers/news"
	"prosnav.com/wxserver/controllers/sync"
	"prosnav.com/wxserver/midwares"
	"prosnav.com/wxserver/controllers/files"
	"prosnav.com/wxserver/controllers/postgres"
	"prosnav.com/wxserver/controllers/activity"
	"prosnav.com/wxserver/controllers/assessment"
	"prosnav.com/wxserver/controllers/appointment"
	"prosnav.com/wxserver/controllers/longhua"
)

func RegisterHandlers(group *gin.RouterGroup) {
    errGroup := group.Group("/error")
    errGroup.Static("", "/var/web/static/publish/error")
    ckGroup := group.Group("/signin", midwares.Auth())
    ckGroup.Static("", "/var/web/static/signin")

    loginGroup := group.Group("/login")
    loginGroup.GET("/act", activity.Authorize)
    loginGroup.GET("", oauth.Authorize)
    loginGroup.POST("/web", oauth.Login)

    oauthGroup := group.Group("/oauth")
    oauthGroup.GET("", oauth.InitAccount)
	//oauthGroup.GET("/", oauth.Authorizion)
	//oauthGroup.GET("/mp", oauth.OauthMP)

    newsGroup := group.Group("/news", midwares.Auth())
    newsGroup.GET("/", news.QueryNewsByUseridWithPage)
    newsGroup.POST("/", news.QueryNewsWithPage)
    newsGroup.GET("/:newsid", news.ExpandNews)
    newsGroup.POST("/delete/:newsid", news.DeleteNewsById)
    newsGroup.POST("/update/:newsid", news.UpdateNews)
    newsGroup.POST("/news", news.InsertNews)

    filesGroup := group.Group("/files", midwares.Auth(), midwares.FileAccess())
    filesGroup.GET("/:tagCode/:newsdir/:filename", files.FileServe)
    filesGroup.POST("/:tagCode/:newsdir", files.Upload)
    filesGroup.POST("/:tagCode", files.Mkdir)
    filesGroup.GET("/:tagCode/:newsdir", files.Ls)


    syncGroup := group.Group("/sync")
    syncGroup.POST("/users", sync.SyncUsers)
    syncGroup.POST("/products", sync.SyncProducts)

    postgresGroup := group.Group("/postgres")
    postgresGroup.POST("/insert", postgres.Insert);
    postgresGroup.POST("/update", postgres.Update);

    actGroup := group.Group("/act")
    actGroup.GET("/meta", midwares.Auth(), activity.InitData)
    actGroup.GET("/activity", midwares.Auth(), activity.QueryActivity)
    actGroup.POST("/activity", activity.AddActivity)
    actGroup.POST("/signin", midwares.Auth(), activity.SignIn)

    assessmentGroup := group.Group("/assessment")
    assessmentGroup.Static("/static", "/var/web/static/assessment")
    assessmentGroup.GET("/paper", assessment.QueryPapper)
    assessmentGroup.POST("", assessment.Assess)
	assessmentGroup.GET("/customers", assessment.CustomersInfo)

    appointmentGroup := group.Group("/appointment")
    appointmentGroup.POST("", appointment.Appointment)

	longhuaGroup := group.Group("/longhua")
	longhuaGroup.POST("", longhua.Reserve)
	longhuaGroup.GET("", longhua.GetReservation)
	longhuaGroup.GET("/time", longhua.GetTime)
	longhuaGroup.GET("/info", longhua.GetAllReservation)

}
