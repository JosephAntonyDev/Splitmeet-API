package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
)

type LoginUserController struct {
	useCase *app.LoginUser
}

func NewLoginUserController(useCase *app.LoginUser) *LoginUserController {
	return &LoginUserController{useCase: useCase}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *LoginUserController) Handle(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	token, err := ctrl.useCase.Execute(req.Email, req.Password)

	if err != nil {
		if err.Error() == "credenciales inválidas" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"token":   token,
	})
}