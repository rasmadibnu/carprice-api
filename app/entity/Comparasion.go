package entity

import (
	"time"

	"gorm.io/gorm"
)

type Comparasion struct {
	ID         int            `gorm:"column:id;primary_key" json:"id"`
	CarsID     int            `gorm:"column:cars_id" json:"cars_id"`
	Title      string         `gorm:"column:title" json:"title"`
	Price      int            `gorm:"column:price" json:"price"`
	SourceName string         `gorm:"column:source_name" json:"source_name"`
	SourceLink string         `gorm:"column:source_link" json:"source_link"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdateAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Comparasion) TableName() string {
	return "tr_comparasion"
}
