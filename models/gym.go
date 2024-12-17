package models

import "gorm.io/gorm"

type Gym struct {
	gorm.Model
	Business_name string `json:"business_name" gorm:"size:255;not null"`
	Address       string `json:"address" gorm:"size:255;not null"`
	City          string `json:"city" gorm:"size:100;not null"`
	Postcode      string `json:"postcode" gorm:"size:20;not null"`
	Phone         string `json:"phone" gorm:"size:15"`
	Email         string `json:"email" gorm:"size:100"`
	Website       string `json:"website" gorm:"size:255"`
	Opening_hours string `json:"opening_hours" gorm:"type:text"`
	Activities    string `json:"activities" gorm:"type:text"`
	Facilities    string `json:"facilities" gorm:"type:text"`
}
