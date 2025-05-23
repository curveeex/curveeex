package handlers

import (
	"curveeex/internal/domain"
	"curveeex/internal/service"
	"curveeex/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtManager  *jwt.JWTManager
}

func NewAuthHandler(authService *service.AuthService, jwtManager *jwt.JWTManager) *AuthHandler {
	return &AuthHandler{authService: authService, jwtManager: jwtManager}
}

// SignUp регистрирует нового пользователя
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req domain.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.authService.SignUp(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

// SignIn аутентифицирует пользователя
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req domain.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	tokens, err := h.authService.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Refresh обновляет токены
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	newTokens, err := h.jwtManager.GenerateTokenPairFromRefresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, newTokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	claims, err := h.jwtManager.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if err := h.authService.Logout(tokenString, int64(claims["exp"].(float64))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Profile data",
	})
}
