package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Validar formato "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use: Bearer <token>"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar el algoritmo de firma (Evita ataques de 'none')
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		// Extraer Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims del token inválidos"})
			return
		}

		// --- AJUSTE CRÍTICO PARA SPLITMEET ---
		// 1. Guardar UserID (Manejando el float64 de JWT)
		if idFloat, ok := claims["user_id"].(float64); ok {
			c.Set("userID", int64(idFloat)) // Convertimos float a int64
		} else {
			// Si por alguna razón el token no tiene ID o viene en otro formato
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token sin user_id válido"})
			return
		}

		// 2. Guardar Email
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		
		// Nota: Eliminé roleID porque tu entidad User actual no tiene roles.
		
		c.Next()
	}
}