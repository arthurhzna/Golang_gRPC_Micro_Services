package handler

import (
	"context"

	"github.com/arthurhzna/Golang_gRPC/internal/utils"
	"github.com/arthurhzna/Golang_gRPC/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer // parent struct from the service package
}

func (ah *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {

	validationErrors, err := utils.CheckValidation(req)
	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("Success"),
	}, nil
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}
