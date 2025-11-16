syntax = "proto3";

package auth;

option go_package = "github.com/arthurhzna/Golang_gRPC/proto/auth";

service AuthService {
	rpc Register(RegisterRequest) returns (RegisterResponse);
}

message RegisterRequest {
	string full_name = 1;
	string email = 2;
	string password = 3;
	string passwrod_confirmation = 4;
}

