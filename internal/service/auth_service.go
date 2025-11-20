package service

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/repository"
	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/auth"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
}

func NewAuthService(authRepository repository.IAuthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}

func (as *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {

	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exists"),
		}, nil
	}

	return &auth.RegisterResponse{}, nil
}
