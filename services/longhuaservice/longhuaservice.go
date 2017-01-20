package longhuaservice

import (
	"prosnav.com/wxserver/domain"
	"prosnav.com/wxserver/utils"
	log "github.com/Sirupsen/logrus"
)

const (
	ALL_RESERVATIONS = `select mobile, name, doctor, time from longhua where mobile = ?`
	ALL_TIMEs = `select time from longhua where doctor = ?`
	ALL_INFO = `select mobile, name, doctor, time from longhua`
)

func CreateReservation(reservation *domain.Reservation) error {
	_, err := utils.Engine.Omit("status").Insert(reservation); if err != nil{
		return err
	}

	return nil
}

func GetReservation(mobile string) (reservations []*domain.Reservation, err error) {
	err = utils.Engine.Sql(ALL_RESERVATIONS, mobile).Find(&reservations); if err != nil {
		log.Error(0, "select reservations info error %v", err)
		return nil, err
	}
	return reservations, nil
}

func GetTime(doctor string) (reservations []*domain.Reservation, err error) {
	err = utils.Engine.Sql(ALL_TIMEs, doctor).Find(&reservations); if err != nil {
		log.Error(0, "select reservations info error %v", err)
		return nil, err
	}
	return reservations, nil
}

func GetAll() (reservations []*domain.ReservationList, err error) {
	err = utils.Engine.Sql(ALL_INFO).Find(&reservations); if err != nil {
		log.Error(0, "select reservations info error %v", err)
		return nil, err
	}
	return reservations, nil
}