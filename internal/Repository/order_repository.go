package repository

import (
	"context"
	"database/sql"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
)

type IOrderRepository interface {
	GetNumbering(ctx context.Context, module string) (*entity.Numbering, error)
	CreateOrder(ctx context.Context, order *entity.Order) error
	UpdateNumbering(ctx context.Context, numbering *entity.Numbering) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (or *orderRepository) GetNumbering(ctx context.Context, module string) (*entity.Numbering, error) {
	row := or.db.QueryRowContext(
		ctx,
		"SELECT module, number FROM numbering WHERE module = $1", module)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var numbering entity.Numbering
	err := row.Scan(&numbering.Module, &numbering.Number)
	if err != nil {
		return nil, err
	}
	return &numbering, nil
}

func (or *orderRepository) CreateOrder(ctx context.Context, order *entity.Order) error {
	_, err := or.db.ExecContext(
		ctx,
		`INSERT INTO "order" (id, number, user_id, order_status_code, user_full_name, address, phone_number, notes, total, expired_at, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		order.Id, order.Number, order.UserId, order.OrderStatusCode, order.UserFullName, order.Address, order.PhoneNumber, order.Notes, order.Total, order.ExpiredAt, order.CreatedAt, order.CreatedBy, order.UpdatedAt, order.UpdatedBy, order.DeletedAt, order.DeletedBy, order.IsDeleted)
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) UpdateNumbering(ctx context.Context, numbering *entity.Numbering) error {
	_, err := or.db.ExecContext(
		ctx,
		"UPDATE numbering SET number = $1 WHERE module = $2", numbering.Number, numbering.Module)
	if err != nil {
		return err
	}
	return nil
}
