package promotion

type Request struct {
	Code      string  `json:"code"  validate:"required"`
	Discount  float64 `json:"discount"  validate:"required"`
	ProductID uint    `json:"product_id"`
}

type CreateRequest struct {
	Request
}

type UpdateRequest struct {
	Request
}
