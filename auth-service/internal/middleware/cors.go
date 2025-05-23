package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Список разрешенных доменов (добавьте свои продакшен-домены при необходимости)
		allowedOrigins := map[string]bool{
			"http://localhost:5173":              true, // Vite dev server
			"http://127.0.0.1:5173":              true, // Альтернативный адрес
			"http://localhost:3000":              true, // Create-React-App по умолчанию
			"https://your-production-domain.com": true,
		}

		// Динамически разрешаем запросы только из доверенных источников
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Обязательные для аутентификации заголовки
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization") // Для доступа к токену
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Authorization, Accept, X-Requested-With, X-CSRF-Token")

		// Разрешаем все основные методы
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE, PATCH")

		// Короткий путь для preflight-запросов
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
