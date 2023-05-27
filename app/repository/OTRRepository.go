package repository

import (
	"api-cariprice/app/entity"
	"api-cariprice/config"
)

type OTRRepository struct {
	config config.Database
}

func NewOTRRepository(database config.Database) OTRRepository {
	return OTRRepository{
		config: database,
	}
}

// @Summary : Get otr
// @Description : -
// @Author : rasmadibnu
func (r *OTRRepository) Search(q string) ([]entity.Cars, error) {
	var otr []entity.Cars

	err := r.config.DB.Preload("Comparasion").Where("title LIKE ?", "%"+q+"%").Find(&otr).Error

	if err != nil {
		return otr, err
	}

	return otr, nil
}

// @Summary : Get otr
// @Description : -
// @Author : rasmadibnu
func (r *OTRRepository) FindAll() ([]entity.Cars, error) {
	var otr []entity.Cars

	err := r.config.DB.Preload("Comparasion").Order("id DESC").Find(&otr).Error

	if err != nil {
		return otr, err
	}

	return otr, nil
}

// @Summary : Insert otr
// @Description : Insert otr to database
// @Author : rasmadibnu
func (r *OTRRepository) Insert(otr entity.Cars) (entity.Cars, error) {
	err := r.config.DB.Debug().Preload("Comparasion").Create(&otr).Error

	if err != nil {
		return otr, err
	}

	return otr, nil
}

// @Summary : Get status
// @Description : find status by ID
// @Author : rasmadibnu
func (r *OTRRepository) FindById(ID int) (entity.Cars, error) {
	var otr entity.Cars

	err := r.config.DB.Where("id = ?", ID).First(&otr).Error

	if err != nil {
		return otr, err
	}

	return otr, nil
}

// // @Summary : Update status
// // @Description : Update status by ID
// // @Author : rasmadibnu
// func (r *StatusRepository) Update(status entity.Status, ID int) (entity.Status, error) {
// 	err := r.config.DB.Where("id = ?", ID).Updates(&status).Error

// 	if err != nil {
// 		return status, err
// 	}

// 	return status, nil
// }

// @Summary : Delete otr
// @Description : Delete otr temporary
// @Author : rasmadibnu
func (r *OTRRepository) Delete(ID int) (bool, error) {
	var otr entity.Cars

	err := r.config.DB.Debug().Where("id = ?", ID).Delete(&otr).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
