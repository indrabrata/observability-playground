package model

type ProductResponse struct {
	Id       int64   `json:"id" example:"1"`
	Name     string  `json:"name" example:"Product A"`
	Quantity int64   `json:"quantity" example:"10"`
	Price    float64 `json:"price" example:"10.99"`
}
