package routes

import (
	"github.com/JosephAntonyDev/splitmeet-api/internal/middleware"
	"github.com/JosephAntonyDev/splitmeet-api/internal/product/infra/controllers"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(
	r *gin.Engine,
	getProductCtrl *controllers.GetProductController,
	getProductsByCategoryCtrl *controllers.GetProductsByCategoryController,
	searchProductsCtrl *controllers.SearchProductsController,
	createCustomProductCtrl *controllers.CreateCustomProductController,
	jwtSecret string,
) {
	g := r.Group("products")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	{
		g.GET("/:id", getProductCtrl.Handle)
		g.GET("/category/:categoryId", getProductsByCategoryCtrl.Handle)
		g.GET("/search", searchProductsCtrl.Handle)
		g.POST("", createCustomProductCtrl.Handle)
	}
}
