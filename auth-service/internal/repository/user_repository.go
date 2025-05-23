package repository

import (
	"curveeex/internal/domain"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Создание пользователя
func (r *UserRepository) Create(user *domain.User) error {
	// УДАЛИТЬ хеширование здесь - пароль уже хеширован в сервисе!
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, user.Email, user.Password).Scan(&user.ID)
}

// Поиск пользователя по email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
