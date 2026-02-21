package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/indrabrata/observability-playground/model"
	"github.com/indrabrata/observability-playground/repository"
	"go.opentelemetry.io/otel/trace"
)

type ProductService struct {
	repository *repository.Queries
	trace      trace.Tracer
}

func New(repository *repository.Queries, trace trace.Tracer) *ProductService {
	return &ProductService{
		repository: repository,
		trace:      trace,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, request model.ProductRequest) (model.ProductResponse, error) {
	ctx, span := s.trace.Start(ctx, "Service.CreateProduct")
	defer span.End()

	product := repository.CreateProductParams{
		Name:      request.Name,
		Quantity:  request.Quantity,
		Price:     request.Price,
		CreatedAt: time.Now(),
	}

	data, err := s.repository.CreateProduct(ctx, product)
	if err != nil {
		return model.ProductResponse{}, err
	}

	response := model.ProductResponse{
		Id:       data.ID,
		Name:     data.Name,
		Quantity: data.Quantity,
		Price:    data.Price,
	}

	return response, nil
}

func (s *ProductService) GetProducts(ctx context.Context) ([]model.ProductResponse, error) {
	data, err := s.repository.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ProductResponse, 0)
	for _, product := range data {
		response := model.ProductResponse{
			Id:       product.ID,
			Name:     product.Name,
			Quantity: product.Quantity,
			Price:    product.Price,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (model.ProductResponse, error) {
	data, err := s.repository.GetProduct(ctx, id)
	if err != nil {
		return model.ProductResponse{}, err
	}

	response := model.ProductResponse{
		Id:       data.ID,
		Name:     data.Name,
		Quantity: data.Quantity,
		Price:    data.Price,
	}

	return response, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int64, request model.ProductRequest) (model.ProductResponse, error) {
	product := repository.UpdateProductParams{
		ID:        id,
		Name:      request.Name,
		Quantity:  request.Quantity,
		Price:     request.Price,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err := s.repository.UpdateProduct(ctx, product)
	if err != nil {
		return model.ProductResponse{}, err
	}

	response := model.ProductResponse{
		Id:       product.ID,
		Name:     product.Name,
		Quantity: product.Quantity,
		Price:    product.Price,
	}

	return response, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	err := s.repository.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
