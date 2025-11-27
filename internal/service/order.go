package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/order"
	"github.com/google/uuid"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error)
}

type orderService struct {
	orderRepository   repository.IOrderRepository
	productRepository repository.IProductRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, productRepository repository.IProductRepository) IOrderService {
	return &orderService{
		orderRepository:   orderRepository,
		productRepository: productRepository,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	numbering, err := os.orderRepository.GetNumbering(ctx, "order")
	if err != nil {
		return nil, err
	}

	var productIds = make([]string, len(req.Products))
	for i, product := range req.Products {
		productIds[i] = product.Id
	}

	products, err := os.productRepository.GetProductsByIds(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]*entity.Product)
	for i := range products {
		productMap[products[i].Id] = products[i]
	}

	var total float64 = 0
	for _, p := range req.Products {
		total += productMap[p.Id].Price * float64(p.Quantity)
	}

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)
	orderEntity := entity.Order{
		Id:              uuid.NewString(),
		Number:          fmt.Sprintf("ORD-%d%08d", now.Year(), numbering.Number),
		UserId:          claims.Subject,
		OrderStatusCode: entity.OrderStatusCodeUnpaid,
		UserFullName:    claims.FullName,
		Address:         req.Address,
		PhoneNumber:     req.PhoneNumber,
		Notes:           &req.Notes,
		Total:           total,
		ExpiredAt:       &expiredAt,
		CreatedAt:       now,
		CreatedBy:       claims.FullName,
	}
	err = os.orderRepository.CreateOrder(ctx, &orderEntity)
	if err != nil {
		return nil, err
	}

	numbering.Number++

	err = os.orderRepository.UpdateNumbering(ctx, numbering)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		Base: utils.SuccessResponse("Order created successfully"),
		Id:   orderEntity.Id,
	}, nil
}
