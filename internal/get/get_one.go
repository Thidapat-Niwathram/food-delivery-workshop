package get

type GetOne[T any] struct {
	ID T `json:"-" path:"id" form:"id" query:"id"`
}

func (f *GetOne[T]) GetID() T {
	return f.ID
}