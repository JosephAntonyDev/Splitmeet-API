package controllers

import (
	"net/http"
	"strconv"

	"github.com/JosephAntonyDev/splitmeet-api/internal/user/app"
	"github.com/gin-gonic/gin"
)

type SearchUsersController struct {
	useCase *app.SearchUsers
}

func NewSearchUsersController(useCase *app.SearchUsers) *SearchUsersController {
	return &SearchUsersController{useCase: useCase}
}

func (c *SearchUsersController) Handle(ctx *gin.Context) {
	query := ctx.Query("username")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username query parameter is required"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	users, err := c.useCase.Execute(query, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if users == nil {
		users = []app.UserSearchResult{}
	}

	ctx.JSON(http.StatusOK, users)
}
