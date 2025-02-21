package models

import (
	"gorm.io/gorm"
)

type Cart struct { // ตะกร้าสินค้า
	gorm.Model
	UserID      uint       `json:"user_id"`
	User        *User      `json:"-" gorm:"foreignKey:UserID"`
	PromotionID *uint      `json:"promotion_id"`
	Promotion   *Promotion `json:"promotion" gorm:"foreignKey:PromotionID"`
	CartItems   []*CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
	SubTotal    float64    `json:"sub_total" gorm:"-"` // รวม CartItem.Price ของ CartItem
	Total       float64    `json:"total" gorm:"-"` // รวมทั้งหมด (หลังหักส่วนลด)
	Discount    float64    `json:"discount" gorm:"-"` //ผลรวมของ Promotion.Discount ของแต่ละ Product
}
