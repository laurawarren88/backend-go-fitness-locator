package models

import "gorm.io/gorm"

type Activities struct {
	gorm.Model
	Name          string  `json:"name" form:"name" gorm:"size:255;not null"`
	Address       string  `json:"address" form:"address" gorm:"size:255;not null"`
	City          string  `json:"city" form:"city" gorm:"size:100;not null"`
	Postcode      string  `json:"postcode" form:"postcode" gorm:"size:20;not null"`
	Description   string  `json:"description" form:"description" gorm:"size:255"`
	Rating        uint    `json:"rating" form:"rating" gorm:"size:1"`
	Phone         string  `json:"phone" form:"phone" gorm:"size:15"`
	Opening_hours string  `json:"opening_hours" form:"opening_hours" gorm:"type:text"`
	TypeID        uint    `json:"type_id" form:"type_id" gorm:"size:255"`
	Type          string  `json:"type" form:"type" gorm:"type:text"`
	Lat           float64 `json:"lat" form:"lat"`
	Lng           float64 `json:"lng" form:"lng"`
	Image         string  `json:"image" form:"image" gorm:"size:255"`
	Website       string  `json:"website" form:"website" gorm:"size:255"`
}

// type Gym struct {
// 	gorm.Model
// 	Phone             string `json:"phone" form:"phone" gorm:"size:15"`
// 	Email             string `json:"email" form:"email" gorm:"size:100"`
// 	Website           string `json:"website" form:"website" gorm:"size:255"`
// 	Opening_hours     string `json:"opening_hours" form:"opening_hours" gorm:"type:text"`
// 	Activities        string `json:"activities" form:"activities" gorm:"type:text"`
// 	Facilities        string `json:"facilities" form:"facilities" gorm:"type:text"`
// 	Logo              string `json:"logo" form:"logo" gorm:"size:255"`
// 	Facilities_images string `json:"facilities_images" form:"facilities_images" gorm:"size:255"`
// }
