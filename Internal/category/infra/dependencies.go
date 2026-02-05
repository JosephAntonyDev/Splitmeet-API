package infra

import (
	"os"

	"github.com/JosephAntonyDev/splitmeet-api/internal/category/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/infra/repository"
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/infra/routes"
	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"github.com/gin-gonic/gin"
)

func SetupDependencies(r *gin.Engine, dbPool *core.Conn_PostgreSQL) {
	// Repository
	categoryRepo := repository.NewCategoryPostgreSQLRepository(dbPool)

	// Use Cases
	getAllCategoriesUseCase := app.NewGetAllCategories(categoryRepo)
	getCategoryUseCase := app.NewGetCategory(categoryRepo)

	// Controllers
	getAllCategoriesController := controllers.NewGetAllCategoriesController(getAllCategoriesUseCase)
	getCategoryController := controllers.NewGetCategoryController(getCategoryUseCase)

	// JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Routes
	routes.SetupCategoryRoutes(r, getAllCategoriesController, getCategoryController, jwtSecret)
}
