package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
)

type ICartRepository interface {
	GetCartByProductAndUserId(ctx context.Context, productId string, userId string) (*entity.Cart, error)
	CreateNewCart(ctx context.Context, cart *entity.Cart) error
	UpdateCart(ctx context.Context, cart *entity.Cart) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) ICartRepository {
	return &cartRepository{db: db}
}

func (cr *cartRepository) GetCartByProductAndUserId(ctx context.Context, productId string, userId string) (*entity.Cart, error) {

	row := cr.db.QueryRowContext(
		ctx,
		"SELECT id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by FROM user_cart WHERE product_id = $1 AND user_id = $2",
		productId,
		userId,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var cart entity.Cart
	err := row.Scan(
		&cart.Id,
		&cart.ProductId,
		&cart.UserId,
		&cart.Quantity,
		&cart.CreatedAt,
		&cart.CreatedBy,
		&cart.UpdatedAt,
		&cart.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (cr *cartRepository) CreateNewCart(ctx context.Context, cart *entity.Cart) error {
	_, err := cr.db.ExecContext(
		ctx,
		`INSERT INTO "user_cart" (id, product_id, user_id, quantity, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		cart.Id,
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.CreatedAt,
		cart.CreatedBy,
		cart.UpdatedAt,
		cart.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) UpdateCart(ctx context.Context, cart *entity.Cart) error {
	_, err := cr.db.ExecContext(
		ctx,
		`UPDATE "user_cart" SET product_id = $1, user_id = $2, quantity = $3, updated_at = $4, updated_by = $5 WHERE id = $6`,
		cart.ProductId,
		cart.UserId,
		cart.Quantity,
		cart.UpdatedAt,
		cart.UpdatedBy,
		cart.Id,
	)

	if err != nil {
		return err
	}
	return nil
}
