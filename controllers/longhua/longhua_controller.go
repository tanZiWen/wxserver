//龙华医院膏方门诊预约
package longhua

import (
	"github.com/gin-gonic/gin"
	log "github.com/Sirupsen/logrus"
	"prosnav.com/wxserver/domain"
	"prosnav.com/wxserver/modules/idg"
	"prosnav.com/wxserver/services/longhuaservice"
	"strings"
)

type ReservationForm struct {
	Mobile string `form:"mobile" json:"mobile"`
	Name   string `form:"name" json:"name"`
	Doctor string `form:"doctor" json:"doctor"`
	Time   string  `form:"time" json:"time"`
}

func Reserve(c *gin.Context) {
	var form ReservationForm
	if err := c.Bind(&form); err != nil {
		c.JSON(500, nil)
		return
	}

	log.Info("reservation form: ", form)
	var time []string
	time = strings.Split(form.Time, ",")
	for i:=0; i < len(time); i++ {
		reservation := new(domain.Reservation)

		userId, err := idg.Id(); if err != nil {
			return
		}

		reservation.Id = userId
		reservation.Mobile = form.Mobile
		reservation.Name = form.Name
		reservation.Doctor = form.Doctor
		log.Debug("reservation time: %v", time[i])
		reservation.Time = time[i]
		err = longhuaservice.CreateReservation(reservation); if err != nil {
			panic(err)
			return
		}
	}

	c.JSON(200, domain.OK)
}

func GetReservation(c *gin.Context)  {
	var reservation []*domain.Reservation
	var err error
	mobile := c.Query("mobile")
	log.Info("mobile", mobile)
	reservation, err = longhuaservice.GetReservation(mobile); if err != nil {
		log.Error(4, "failed to get reservation", err)
		panic(err)
	}
	c.JSON(200, reservation)
}

func GetTime(c *gin.Context)  {
	var reservation []*domain.Reservation
	var err error
	doctor := c.Query("doctor")
	log.Info("doctor", doctor)
	reservation, err = longhuaservice.GetTime(doctor); if err != nil {
		log.Error(4, "failed to get reservation", err)
		panic(err)
	}
	c.JSON(200, reservation)
}

func GetAllReservation(c *gin.Context) {
	var reservation []*domain.ReservationList
	var err error
	reservation, err = longhuaservice.GetAll(); if err != nil {
		log.Error(4, "failed to get reservation", err)
		panic(err)
	}
	c.JSON(200, reservation)
}

