package postgres

import (
	"github.com/gin-gonic/gin"
	"prosnav.com/wxserver/services/userservice"
	"prosnav.com/wxserver/modules/log"
)

type Error struct  {
	Userid string
	Error error
}

func Insert(c *gin.Context) {
	arr := []string{
		"18717876582",
		"tanyuan",
	}
	tag := "10023"
	errMap := make(map[string]error, 0)
	userservice.InsertUser(arr, []string{tag}, &errMap);
	for k, v := range errMap {
		log.Debug("%s Error is:%v\n", k, v)
	}
	c.JSON(200, nil)
}

func Update(c *gin.Context) {
	arr :=[]string{
		"13661538136"}

	tag := "10023"

	err := userservice.UpdateUser(arr, tag); if err != nil {
		panic(err)
	}

	c.JSON(200, nil)
}

