package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/indrabrata/observability-playground/docs"
	"github.com/indrabrata/observability-playground/handler"
	"github.com/indrabrata/observability-playground/infrastructure"
	"github.com/indrabrata/observability-playground/middleware"
	"github.com/indrabrata/observability-playground/repository"
	"github.com/indrabrata/observability-playground/service"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

// Span is a unit that records particular operation within certain time window.

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  contact@ndrz.dev
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		zap.L().Fatal("failed to load environment variables")
	}

	infrastructure.NewZapLog(ctx)

	db := infrastructure.SqlLite3DBConnect(ctx)
	defer db.Close()

	metric := infrastructure.NewPrometheusMetric(ctx)
	trace := infrastructure.NewOpenTelemetryTrace(ctx)

	router := chi.NewRouter()
	router.Use(middleware.RequestIdMiddleware)
	router.Use(middleware.MetricsMiddleware)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	productRepository := repository.New(db)
	productSService := service.New(productRepository, trace.Tracer("Product.Service"))
	productHandler := handler.New(productSService, trace.Tracer("Product.Handler"))

	router.Post("/products", productHandler.CreateProduct)
	router.Get("/products", productHandler.GetProducts)
	router.Get("/products/{id}", productHandler.GetProduct)
	router.Put("/products/{id}", productHandler.UpdateProduct)
	router.Delete("/products/{id}", productHandler.DeleteProduct)
	router.Get("/metrics", promhttp.HandlerFor(metric, promhttp.HandlerOpts{}).ServeHTTP)

	port := os.Getenv("PORT")
	zap.L().Info("Starting server on port " + port)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	docs.SwaggerInfo.Title = "Observability Playground API"
	docs.SwaggerInfo.Description = "This is a sample server for observability."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	http.ListenAndServe(":"+port, router)
}
