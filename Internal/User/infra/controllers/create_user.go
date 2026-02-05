package controllers

import (
	"net/http"
	"strings"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/domain/entities"
	"github.com/gin-gonic/gin"
)

type CreateUserController struct {
	useCase *app.CreateUser
}

func NewCreateUserController(useCase *app.CreateUser) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`
}

// El método que Gin va a ejecutar
func (ctrl *CreateUserController) Handle(c *gin.Context) {
	var req createUserRequest

	// 1. Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// 2. Mapear DTO -> Entidad de Dominio
	user := entities.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}

	// 3. Ejecutar UseCase (Pasando el puntero &user)
	err := ctrl.useCase.Execute(&user)

	// 4. Manejo de Errores
	if err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Respuesta Exitosa (201 Created)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Usuario creado exitosamente",
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
