package models

import "gorm.io/gorm"

type Place struct {
	gorm.Model
	Name            string  `json:"name" form:"name" gorm:"size:255" binding:"required"`
	Vicinity        string  `json:"vicinity" form:"vicinity" gorm:"size:255"`
	City            string  `json:"city" form:"city" gorm:"size:100"`
	Postcode        string  `json:"postcode" form:"postcode" gorm:"size:20"`
	Phone           string  `json:"phone" form:"phone" gorm:"size:15" binding:"required"`
	Email           string  `json:"email" form:"email" gorm:"size:100"`
	Website         string  `json:"website" form:"website" gorm:"size:255"`
	OpeningHours    string  `json:"opening_hours" form:"opening_hours" gorm:"type:text"`
	TypeID          uint    `json:"type_id" form:"type_id"`
	Type            string  `json:"type" form:"type" gorm:"type:text"`
	Description     string  `json:"description" form:"description" gorm:"size:255" binding:"required"`
	Photo           string  `json:"photo" form:"photo" gorm:"size:255"`
	FacilitiesImage string  `json:"facilities_image" form:"facilities_image" gorm:"size:255"`
	Rating          uint    `json:"rating" form:"rating"`
	Latitude        float64 `json:"latitude" form:"latitude"`
	Longitude       float64 `json:"longitude" form:"longitude"`
}
