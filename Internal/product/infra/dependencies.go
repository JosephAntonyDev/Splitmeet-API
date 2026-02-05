package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/infra/routes"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Repository
	productRepo := repository.NewProductPostgreSQLRepository(dbPool)

	// Use Cases
	getProductUseCase := app.NewGetProduct(productRepo)
	getProductsByCategoryUseCase := app.NewGetProductsByCategory(productRepo)
	searchProductsUseCase := app.NewSearchProducts(productRepo)
	createCustomProductUseCase := app.NewCreateCustomProduct(productRepo)

	// Controllers
	getProductController := controllers.NewGetProductController(getProductUseCase)
	getProductsByCategoryController := controllers.NewGetProductsByCategoryController(getProductsByCategoryUseCase)
	searchProductsController := controllers.NewSearchProductsController(searchProductsUseCase)
	createCustomProductController := controllers.NewCreateCustomProductController(createCustomProductUseCase)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupProductRoutes(
		r,
		getProductController,
		getProductsByCategoryController,
		searchProductsController,
		createCustomProductController,
		jwtSecret,
	)
}
