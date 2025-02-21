package cart

type CartItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

type CreateRequest  struct {
	CartItemRequests  []CartItemRequest `json:"cart_items"`
}

type UpdateRequest struct {
	CartItemRequests  []CartItemRequest `json:"cart_items"`
}

type PromotionRequest struct {
	PromotionCode string `json:"promotion_code" validate:"required"`
}