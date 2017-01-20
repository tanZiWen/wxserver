package appointmentservice

import (
	"prosnav.com/wxserver/domain"
	"prosnav.com/wxserver/utils"
)

func MakeAppointment(appoint *domain.Appoint) error {
	_, err := utils.Engine.Omit("status").Insert(appoint); if err != nil{
		return err
	}

	return nil
}

func GetMobile(appoint *domain.Appoint) (bool, error) {
	return utils.Engine.Get(appoint)
}
