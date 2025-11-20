package repository

import (
	"context"
	"database/sql"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{db: db}
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, FROM user WHERE email = $1 AND is_delated IS false", email)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
