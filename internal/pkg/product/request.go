package product

type Request struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
}

type CreateRequest struct {
	Request
}

type UpdateRequest struct {
	Request
}
