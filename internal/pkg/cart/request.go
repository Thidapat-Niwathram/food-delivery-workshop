package cart

type CartItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

type CreateRequest struct {
	UserID           uint              `json:"-"`
	CartItemRequests []CartItemRequest `json:"cart_items"`
}

type UpdateRequest struct {
	UserID           uint              `json:"-"`
	CartItemRequests []CartItemRequest `json:"cart_items"`
}

type PromotionRequest struct {
	UserID        uint   `json:"-"`
	PromotionCode string `json:"promotion_code" validate:"required"`
}

type RemoveItemRequest struct {
	UserID    uint `json:"-"`
	ProductID uint `json:"product_id" validate:"required"`
}

type GetAllRequests struct {
	UserID uint `json:"-"`
}