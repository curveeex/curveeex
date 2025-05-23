package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// JWTManager структура должна быть объявлена в jwt.go
type JWTManager struct {
	secretKey     string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

func (m *JWTManager) AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		// Проверка в Redis
		if rdb != nil {
			val, err := rdb.Get(context.Background(), "bl_"+tokenString).Result()
			if err == nil && val == "1" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
				return
			}
		}

		claims, err := m.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims["user_id"])
		c.Next()
	}
}

// VerifyToken должен быть объявлен в jwt.go
func (m *JWTManager) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
