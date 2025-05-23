package service

import (
	"context"
	"curveeex/internal/domain"
	"curveeex/internal/repository"
	"curveeex/pkg/jwt"
	"errors"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	jwtManager  *jwt.JWTManager
	redisClient *redis.Client
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwt.JWTManager, rdb *redis.Client) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtManager:  jwtManager,
		redisClient: rdb,
	}
}

func (s *AuthService) SignUp(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	log.Printf("Original password: %s", user.Password)
	log.Printf("Hashed password: %s", string(hashedPassword))

	user.Password = string(hashedPassword)
	return s.userRepo.Create(user)
}

func (s *AuthService) SignIn(email, password string) (*jwt.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("User not found: %v", err)
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)
		return nil, errors.New("invalid password")
	}

	return s.jwtManager.GenerateTokenPair(user.ID)
}

func (s *AuthService) Logout(token string, exp int64) error {
	if s.redisClient == nil {
		return nil
	}

	ttl := time.Unix(exp, 0).Sub(time.Now())
	if ttl > 0 {
		return s.redisClient.SetNX(context.Background(), "bl_"+token, "1", ttl).Err()
	}
	return nil
}
