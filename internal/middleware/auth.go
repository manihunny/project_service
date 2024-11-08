package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.fast-go.ru/fast-go-team/auth/pkg/authrpc"
	"gitlab.fast-go.ru/fast-go-team/project/config"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization token required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("Bearer token required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		authClient := authrpc.NewAuthClient(config.NewAppConfig().AuthServiceGRPCAddress)
		tokenReq := authrpc.VerifyTokenRequest{Token: tokenString}

		claims, err := authClient.VerifyToken(c, &tokenReq)
		if err != nil {
			log.Println("Invalid token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Printf("UserID extracted from token: %v", claims.GetUserId())
		c.Set("userID", claims.GetUserId())
		c.Set("role", claims.GetServiceRole())
		c.Next()
	}
}
