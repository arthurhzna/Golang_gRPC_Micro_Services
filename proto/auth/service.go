syntax = "proto3";

package auth;

option go_package = "github.com/arthurhzna/Golang_gRPC/pb/service";

service AuthService {
	rpc Register(RegisterRequest) returns (RegisterResponse);
}

message RegisterRequest {
	string full_name = 1;
	string email = 2;
	string password = 3;
	string password_confirmation = 4;
}

message RegisterResponse {
	string message = 1;
}
