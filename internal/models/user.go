package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Password       string `json:"password" validate:"required"`
	Phone          string `json:"phone" `
	IDCard         string `json:"id_card"`
	Address        string `json:"address"`
	AddressDetails string `json:"address_details"`
	Cart           []Cart `json:"-" gorm:"foreignKey:UserID"`
}

