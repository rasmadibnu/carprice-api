package entity

type Comparasion struct {
	ID         int    `gorm:"column:id;primary_key" json:"id"`
	CarsID     int    `gorm:"column:cars_id" json:"cars_id"`
	Title      string `gorm:"column:title" json:"title"`
	Price      int    `gorm:"column:price" json:"price"`
	SourceName string `gorm:"column:source_name" json:"source_name"`
	SourceLink string `gorm:"column:source_link" json:"source_link"`
}

func (Comparasion) TableName() string {
	return "tr_comparasion"
}
