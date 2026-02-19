package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/indrabrata/observability-playground/docs"
	"github.com/indrabrata/observability-playground/handler"
	"github.com/indrabrata/observability-playground/infrastructure"
	"github.com/indrabrata/observability-playground/repository"
	"github.com/indrabrata/observability-playground/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  contact@ndrz.dev
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load environment variables")
	}

	port := os.Getenv("PORT")

	docs.SwaggerInfo.Title = "Observability Playground API"
	docs.SwaggerInfo.Description = "This is a sample server for observability."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	db := infrastructure.SqlLite3Connect()
	defer db.Close()

	productRepository := repository.New(db)
	productSService := service.New(productRepository)
	productHandler := handler.New(productSService)

	router.Post("/products", productHandler.CreateProduct)
	router.Get("/products", productHandler.GetProducts)
	router.Get("/products/{id}", productHandler.GetProduct)
	router.Put("/products/{id}", productHandler.UpdateProduct)
	router.Delete("/products/{id}", productHandler.DeleteProduct)

	log.Info().Msg("Starting server on port " + port)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+port+"/swagger/doc.json"),
	))

	http.ListenAndServe(":"+port, router)
}
