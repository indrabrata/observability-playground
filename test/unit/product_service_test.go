package unit

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/indrabrata/observability-playground/model"
	"github.com/indrabrata/observability-playground/repository"
	productrepository "github.com/indrabrata/observability-playground/repository/product"
	"github.com/indrabrata/observability-playground/service"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestCreateProduct(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	tracer := noop.NewTracerProvider().Tracer("test")

	productRepository := productrepository.New(db)
	productService := service.New(repository.NewBaseRepository(db, productRepository), tracer)

	ctx := context.WithValue(context.Background(), "requestId", "test-123")

	request := model.ProductRequest{Name: "Test Product", Quantity: 10, Price: 100.0}
	result, err := productService.CreateProduct(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, int64(10), result.Quantity)
	assert.Equal(t, 100.0, result.Price)
}
