package service

import (
	"api-cariprice/app/entity"
	"api-cariprice/app/repository"
	"time"
)

type OTRServices struct {
	repository repository.OTRRepository
}

func NewOTRServices(r repository.OTRRepository) OTRServices {
	return OTRServices{
		repository: r,
	}
}

// @Summary : Search otr
// @Description : Search otr from repository
// @Author : rasmadibnu
func (s *OTRServices) Search(q string) ([]entity.Cars, int, error) {
	status, err := s.repository.Search(q)

	if err != nil {
		return status, len(status), err
	}

	return status, len(status), nil
}

// @Summary : List otr
// @Description : Get otr from repository
// @Author : rasmadibnu
func (s *OTRServices) List() ([]entity.Cars, error) {
	status, err := s.repository.FindAll()

	if err != nil {
		return status, err
	}

	return status, nil
}

// @Summary : Insert status
// @Description : insert status to repository
// @Author : rasmadibnu
func (s *OTRServices) Insert(otr entity.Cars) (entity.Cars, error) {
	otr.UpdateAt = time.Now()
	newOTR, err := s.repository.Insert(otr)

	if err != nil {
		return newOTR, err
	}

	return newOTR, nil
}

// @Summary : Find status
// @Description : Find status by id from repository
// @Author : rasmadibnu
func (s *OTRServices) FindById(ID int) (entity.Cars, error) {
	status, err := s.repository.FindById(ID)

	if err != nil {
		return status, err
	}

	return status, nil
}

// // @Summary : Update status
// // @Description : Update status by id from repository
// // @Author : rasmadibnu
// func (s *StatusService) Update(req request.Status, ID int) (entity.Status, error) {
// 	status := entity.Status{
// 		Name:        req.Name,
// 		Color:       req.Color,
// 		Description: req.Description,
// 	}

// 	updateStatus, err := s.repository.Update(status, ID)

// 	if err != nil {
// 		return updateStatus, err
// 	}

// 	return updateStatus, nil
// }

// @Summary : Delete otr
// @Description : Delete otr from repository
// @Author : rasmadibnu
func (s *OTRServices) Delete(ID int) (bool, error) {
	otr, err := s.repository.Delete(ID)

	if err != nil {
		return false, err
	}

	return otr, nil
}
