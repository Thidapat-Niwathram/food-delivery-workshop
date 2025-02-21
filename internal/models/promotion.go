package models

import (
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	Code      string  `json:"code"`
	Discount  float64 `json:"discount"`
	ProductID uint    `json:"product_id"`
	Product   *Product `json:"product" gorm:"foreignKey:ProductID"`
}
