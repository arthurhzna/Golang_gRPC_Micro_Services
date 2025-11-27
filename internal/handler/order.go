package handler

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/order"
)

type orderHandler struct {
	order.UnimplementedOrderServiceServer

	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (oh *orderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &order.CreateOrderResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := oh.orderService.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
