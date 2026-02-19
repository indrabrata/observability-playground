package model

import "errors"

type ProductRequest struct {
	Name     string  `json:"name" example:"Product A"`
	Quantity int64   `json:"quantity" example:"10"`
	Price    float64 `json:"price" example:"10.99"`
}

func (pr *ProductRequest) Validate() error {
	if pr.Name == "" {
		return errors.New("name is required")
	}
	if pr.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if pr.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}
