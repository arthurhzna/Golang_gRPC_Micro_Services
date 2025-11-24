package service

import (
	"context"
	"time"

	"github.com/arthurhzna/Golang_gRPC/internal/entity"
	jwtentity "github.com/arthurhzna/Golang_gRPC/internal/entity/jwt"
	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/cart"
	"github.com/google/uuid"
)

type ICartService interface {
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
}

type cartService struct {
	productRepository repository.IProductRepository
	cartRepository    repository.ICartRepository
}

func NewCartService(productRepository repository.IProductRepository, cartRepository repository.ICartRepository) ICartService {
	return &cartService{
		productRepository: productRepository,
		cartRepository:    cartRepository,
	}
}

func (cs *cartService) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {

	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	productEntity, err := cs.productRepository.GetProductById(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	if productEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	cartEntity, err := cs.cartRepository.GetCartByProductAndUserId(ctx, req.ProductId, claims.Subject)
	if err != nil {
		return nil, err
	}
	if cartEntity == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("Cart not found"),
		}, nil
	}

	if cartEntity != nil {
		now := time.Now()
		cartEntity.Quantity += 1
		cartEntity.UpdatedAt = &now
		cartEntity.UpdatedBy = &claims.Subject

		err = cs.cartRepository.UpdateCart(ctx, cartEntity)
		if err != nil {
			return nil, err
		}

		return &cart.AddProductToCartResponse{
			Base: utils.SuccessResponse("Product added to cart successfully"),
			Id:   cartEntity.Id,
		}, nil
	}

	newCartEntity := entity.Cart{
		Id:        uuid.NewString(),
		UserId:    claims.Subject,
		ProductId: req.ProductId,
		Quantity:  cartEntity.Quantity + 1,
		CreatedAt: time.Now(),
		CreatedBy: claims.FullName,
	}

	err = cs.cartRepository.CreateNewCart(ctx, &newCartEntity)
	if err != nil {
		return nil, err
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse("Product added to cart successfully"),
		Id:   productEntity.Id,
	}, nil
}
