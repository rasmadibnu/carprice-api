package entity

type FilterOption struct {
	Found     int        `json:"found"`
	TotalData int        `json:"total_data"`
	Location  []Location `json:"location"`
}

type Location struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}
