package utils

import (
	"github.com/arthurhzna/Golang_gRPC/pb/common"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
		IsError:    false,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidateError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:     400,
		Message:        "Validation errors",
		IsError:        true,
		ValidateErrors: validationErrors,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}
