package main

import (
	"curveeex/internal/config"
	"curveeex/internal/handlers"
	"curveeex/internal/middleware"
	"curveeex/internal/repository"
	"curveeex/internal/service"
	"curveeex/pkg/jwt"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"log"
	"time"

	// Swagger imports
	_ "curveeex/docs" // Make sure you have this docs package generated
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description REST API для аутентификации пользователей
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	// Подключение к PostgreSQL
	db, err := sql.Open("postgres",
		"host="+cfg.DB.Host+" "+
			"port="+cfg.DB.Port+" "+
			"user="+cfg.DB.User+" "+
			"password="+cfg.DB.Password+" "+
			"dbname="+cfg.DB.Name+" sslmode=disable")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	// Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Инициализация JWT
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret,
		time.Hour*24,   // Access token expire (24h)
		time.Hour*24*7) // Refresh token expire (7d)

	// Инициализация сервисов
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager, rdb)
	authHandler := handlers.NewAuthHandler(authService, jwtManager)

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/signup", authHandler.SignUp)
	router.POST("/signin", authHandler.SignIn)
	router.POST("/refresh", authHandler.Refresh)

	// Protected routes
	authGroup := router.Group("/")
	authGroup.Use(jwtManager.AuthMiddleware(rdb))
	{
		authGroup.GET("/profile", authHandler.GetProfile)
		authGroup.POST("/logout", authHandler.Logout)
	}

	log.Println("Server running on :8080")
	log.Println("Swagger UI at http://localhost:8080/swagger/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server error:", err)
	}
}
