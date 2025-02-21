package models

import (
	"gorm.io/gorm"
)

type CartItem struct { //รายการในตะกร้าสินค้า
	gorm.Model
	CartID     uint     `json:"cart_id"`
	Cart       *Cart    `json:"cart" gorm:"foreignKey:CartID"`
	ProductID  uint     `json:"product_id"`
	Quantity   uint     `json:"quantity"`
	Price      float64  `json:"price" gorm:"-"` // ราคาสินค้า Product.Price
	TotalPrice float64  `json:"total_price" gorm:"-"`
	Product    *Product `json:"product" gorm:"foreignKey:ProductID"`
}

func (ci *CartItem) AfterFind(tx *gorm.DB) (err error) {
	ci.CalculatePrice()
	return nil
}


func (ci *CartItem) CalculatePrice() {
	if ci.Product != nil {
		ci.Price = ci.Product.Price
		ci.TotalPrice = ci.Price * float64(ci.Quantity)
	}
}