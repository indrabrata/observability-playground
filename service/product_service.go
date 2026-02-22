package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/indrabrata/observability-playground/model"
	"github.com/indrabrata/observability-playground/repository"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
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

	zap.L().Debug("create product payload", zap.String("requestId", ctx.Value("requestId").(string)),
		zap.Dict("product",
			zap.String("name", product.Name),
			zap.Int64("quantity", product.Quantity),
			zap.Float64("price", product.Price)))

	data, err := s.repository.CreateProduct(ctx, product)
	if err != nil {
		zap.L().Error("failed to create product", zap.Error(err), zap.String("requestId", ctx.Value("requestId").(string)))
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
	ctx, span := s.trace.Start(ctx, "Service.GetProducts")
	defer span.End()

	data, err := s.repository.GetProducts(ctx)
	if err != nil {
		zap.L().Error("failed to get products", zap.Error(err), zap.String("requestId", ctx.Value("requestId").(string)))
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
	ctx, span := s.trace.Start(ctx, "Service.GetProduct")
	defer span.End()

	data, err := s.repository.GetProduct(ctx, id)
	if err != nil {
		zap.L().Error("failed to get product", zap.Error(err), zap.String("requestId", ctx.Value("requestId").(string)))
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
	ctx, span := s.trace.Start(ctx, "Service.UpdateProduct")
	defer span.End()

	product := repository.UpdateProductParams{
		ID:        id,
		Name:      request.Name,
		Quantity:  request.Quantity,
		Price:     request.Price,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	zap.L().Debug("update product payload", zap.String("requestId", ctx.Value("requestId").(string)),
		zap.Dict("product",
			zap.String("name", product.Name),
			zap.Int64("quantity", product.Quantity),
			zap.Float64("price", product.Price)))

	err := s.repository.UpdateProduct(ctx, product)
	if err != nil {
		zap.L().Error("failed to update product", zap.Error(err), zap.String("requestId", ctx.Value("requestId").(string)))
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
	ctx, span := s.trace.Start(ctx, "Service.DeleteProduct")
	defer span.End()

	err := s.repository.DeleteProduct(ctx, id)
	if err != nil {
		zap.L().Error("failed to delete product", zap.Error(err), zap.String("requestId", ctx.Value("requestId").(string)))
		return err
	}

	return nil
}
