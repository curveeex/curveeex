package jwt

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type JWTManager struct {
	secretKey     string
	accessExpire  time.Duration
	refreshExpire time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

func NewJWTManager(secretKey string, accessExpire, refreshExpire time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

func (m *JWTManager) GenerateTokenPair(userID int) (*TokenPair, error) {
	accessToken, err := m.generateToken(userID, m.accessExpire, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.generateToken(userID, m.refreshExpire, "refresh")
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) generateToken(userID int, expire time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expire).Unix(),
		"type":    tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func (m *JWTManager) GenerateTokenPairFromRefresh(refreshToken string) (*TokenPair, error) {
	claims, err := m.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims["type"] != "refresh" {
		return nil, errors.New("not a refresh token")
	}

	userID := int(claims["user_id"].(float64))
	return m.GenerateTokenPair(userID)
}
