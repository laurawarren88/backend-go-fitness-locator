package models

import "gorm.io/gorm"

type Place struct {
	gorm.Model
	Name            string  `json:"name" form:"name" gorm:"size:255;not null"`
	Vicinity        string  `json:"vicinity" form:"vicinity" gorm:"size:255"`
	City            string  `json:"city" form:"city" gorm:"size:100;not null"`
	Postcode        string  `json:"postcode" form:"postcode" gorm:"size:20;not null"`
	Phone           string  `json:"phone" form:"phone" gorm:"size:15"`
	Email           string  `json:"email" form:"email" gorm:"size:100"`
	Website         string  `json:"website" form:"website" gorm:"size:255"`
	OpeningHours    string  `json:"opening_hours" form:"opening_hours" gorm:"type:text"`
	TypeID          uint    `json:"type_id" form:"type_id" gorm:"not null"`
	Type            string  `json:"type" form:"type" gorm:"type:text"`
	Description     string  `json:"description" form:"description" gorm:"size:255"`
	Photo           string  `json:"photo" form:"photo" gorm:"size:255"`
	FacilitiesImage string  `json:"facilities_image" form:"facilities_image" gorm:"size:255"`
	Rating          uint    `json:"rating" form:"rating"`
	Latitude        float64 `json:"latitude" form:"latitude"`
	Longitude       float64 `json:"longitude" form:"longitude"`
}
