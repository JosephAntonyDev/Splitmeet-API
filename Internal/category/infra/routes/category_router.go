package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/category/infra/controllers"
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(
	r *gin.Engine,
	getAllCategoriesCtrl *controllers.GetAllCategoriesController,
	getCategoryCtrl *controllers.GetCategoryController,
	jwtSecret string,
) {
	g := r.Group("categories")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	{
		g.GET("", getAllCategoriesCtrl.Handle)
		g.GET("/:id", getCategoryCtrl.Handle)
	}
}
