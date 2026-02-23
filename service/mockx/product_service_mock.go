package mockx

import (
	"context"

	"github.com/indrabrata/observability-playground/model"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func NewProductServiceMock() *ProductServiceMock {
	return &ProductServiceMock{}
}

func (m *ProductServiceMock) CreateProduct(ctx context.Context, product model.ProductRequest) (model.ProductResponse, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) GetProducts(ctx context.Context) ([]model.ProductResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) GetProduct(ctx context.Context, id int64) (model.ProductResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) UpdateProduct(ctx context.Context, product model.ProductRequest) (model.ProductResponse, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(model.ProductResponse), args.Error(1)
}

func (m *ProductServiceMock) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
