package entity

import (
	"time"

	"gorm.io/gorm"
)

type Cars struct {
	ID              int            `gorm:"column:id;primary_key" json:"id"`
	Thumbnail       string         `gorm:"column:thumbnail" json:"thumbnail"`
	Title           string         `gorm:"column:title" json:"title"`
	Description     string         `gorm:"column:description" json:"description"`
	Price           int            `gorm:"column:price" json:"price"`
	Brand           string         `gorm:"column:brand" json:"brand"`
	Model           string         `gorm:"column:model" json:"model"`
	Location        string         `gorm:"column:location" json:"location"`
	Detail          string         `gorm:"column:detail" json:"detail"`
	Fuel            string         `gorm:"column:fuel" json:"fuel"`
	IsScraping      bool           `gorm:"-" json:"is_scraping"`
	Status          string         `gorm:"column:status" json:"status"`
	SourceName      string         `gorm:"column:source_name" json:"source_name"`
	SourceLink      string         `gorm:"column:source_link" json:"source_link"`
	Total           int            `gorm:"column:total" json:"total"`
	Average         int            `gorm:"column:average" json:"average"`
	Depreciation    int            `gorm:"column:depreciation" json:"depreciation"`
	Comparasion     []Comparasion  `gorm:"foreignKey:CarsID" json:"comparasion"`
	CreatedAtString string         `gorm:"column:created_at_string" json:"created_at"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at_time"`
	UpdateAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Cars) TableName() string {
	return "ms_cars"
}

func (c *Cars) AfterFind(tx *gorm.DB) (err error) {
	c.IsScraping = false
	return
}
