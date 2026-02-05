package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	categoryInfra "github.com/JosephAntonyDev/splitmeet-api/internal/category/infra"
	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	groupInfra "github.com/JosephAntonyDev/splitmeet-api/internal/group/infra"
	outingInfra "github.com/JosephAntonyDev/splitmeet-api/internal/outing/infra"
	paymentInfra "github.com/JosephAntonyDev/splitmeet-api/internal/payment/infra"
	productInfra "github.com/JosephAntonyDev/splitmeet-api/internal/product/infra"
	userInfra "github.com/JosephAntonyDev/splitmeet-api/internal/user/infra"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Conectar a la Base de Datos
	dbPool, err := core.GetDBPool()
	if err != nil {
		log.Fatalf("Error fatal al conectar con la base de datos: %v", err)
	}
	defer dbPool.DB.Close()

	r := gin.Default()

	r.Use(core.SetupCORS())

	// Inyectar Dependencias
	userInfra.SetupDependencies(r, dbPool)
	categoryInfra.SetupDependencies(r, dbPool)
	productInfra.SetupDependencies(r, dbPool)
	groupInfra.SetupDependencies(r, dbPool)
	outingInfra.SetupDependencies(r, dbPool)
	paymentInfra.SetupDependencies(r, dbPool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor Splitmeet iniciado en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
