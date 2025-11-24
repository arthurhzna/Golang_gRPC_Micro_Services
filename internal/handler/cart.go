package handler

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/service"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/cart"
)

type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

type ICartService interface {
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
}

func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

func (ch *cartHandler) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &cart.AddProductToCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	res, err := ch.cartService.AddProductToCart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
