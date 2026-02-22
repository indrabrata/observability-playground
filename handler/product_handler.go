package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/indrabrata/observability-playground/model"
	"github.com/indrabrata/observability-playground/service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service *service.ProductService
	trace   trace.Tracer
}

func New(service *service.ProductService, trace trace.Tracer) *ProductHandler {
	return &ProductHandler{
		service: service,
		trace:   trace,
	}
}

// @Summary Create product
// @Description Creates a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param CreateProduct body model.ProductRequest true "Product details"
// @Success 200 {object} model.ProductResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("creating product", zap.String("requestId", r.Context().Value("requestId").(string)))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	ctx, span := h.trace.Start(ctx, "Handler.CreateProduct", trace.WithAttributes(attribute.String("requestId", r.Context().Value("requestId").(string))))
	defer span.End()

	var req model.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		zap.L().Error("failed to decode product request", zap.Error(err), zap.String("requestId", r.Context().Value("requestId").(string)))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		zap.L().Error("failed to validate product request", zap.Error(err), zap.String("requestId", r.Context().Value("requestId").(string)))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zap.L().Info("product created", zap.String("requestId", r.Context().Value("requestId").(string)), zap.Int64("productId", product.Id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// @Summary Get products
// @Description Retrieves all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} []model.ProductResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	zap.L().Info("get products", zap.String("requestId", r.Context().Value("requestId").(string)))

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	ctx, span := h.trace.Start(ctx, "Handler.GetProducts", trace.WithAttributes(attribute.String("requestId", r.Context().Value("requestId").(string))))
	defer span.End()

	products, err := h.service.GetProducts(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zap.L().Info("products retrieved", zap.String("requestId", r.Context().Value("requestId").(string)), zap.Int("productCount", len(products)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// @Summary Get product by ID
// @Description Retrieves product based on given ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.ProductResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	ctx, span := h.trace.Start(ctx, "Handler.GetProduct", trace.WithAttributes(attribute.String("requestId", r.Context().Value("requestId").(string))))
	defer span.End()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		zap.L().Error("failed to parse id", zap.String("requestId", ctx.Value("requestId").(string)), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zap.L().Info("getting product", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zap.L().Info("product retrieved", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// @Summary Update product
// @Description Updates a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param UpdateProduct body model.ProductRequest true "Product details"
// @Success 200 {object} model.ProductResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	ctx, span := h.trace.Start(ctx, "Handler.UpdateProduct", trace.WithAttributes(attribute.String("requestId", r.Context().Value("requestId").(string))))
	defer span.End()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zap.L().Info("updating product", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	var req model.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.service.UpdateProduct(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zap.L().Info("product updated", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// @Summary Delete product
// @Description Deletes a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	ctx, span := h.trace.Start(ctx, "Handler.DeleteProduct", trace.WithAttributes(attribute.String("requestId", r.Context().Value("requestId").(string))))
	defer span.End()

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zap.L().Info("deleting product", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	err = h.service.DeleteProduct(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zap.L().Info("product deleted", zap.String("requestId", ctx.Value("requestId").(string)), zap.Int64("productId", id))

	w.WriteHeader(http.StatusNoContent)
}
