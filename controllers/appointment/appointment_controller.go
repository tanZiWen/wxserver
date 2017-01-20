package appointment

import (
	"github.com/gin-gonic/gin"
	log "github.com/Sirupsen/logrus"
	"prosnav.com/wxserver/domain"
	"prosnav.com/wxserver/services/appointmentservice"
	"prosnav.com/wxserver/modules/idg"
)

type AppointForm struct {
	Mobile string   `form:"mobile" json:"mobile"`
	Name   string `form:"name" json:"name"`
	Note   string `form:"note" json"note"`
	Number int `form:"number" json:"number"`
}

func Appointment(c *gin.Context) {
	var form AppointForm
	if err := c.Bind(&form); err != nil {
		c.JSON(500, nil)
		return
	}

	log.Debug("appoint form: %+v", form)
	appoint := new(domain.Appoint)

	appoint.Mobile = form.Mobile
	isExist, err := appointmentservice.GetMobile(appoint); if err != nil {
		log.Error("failed to find user by mobile, %v", err)
		panic(err)
	}

	if isExist {
		c.JSON(200, domain.ISEXIST)
		return
	}

	userId, err := idg.Id(); if err != nil {
		return
	}

	appoint.Id = userId
	appoint.Name = form.Name
	appoint.Note = form.Note
	appoint.Number = form.Number
	err = appointmentservice.MakeAppointment(appoint); if err != nil {
		panic(err)
		return
	}

	c.JSON(200, domain.OK)
}
