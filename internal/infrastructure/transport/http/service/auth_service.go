package service

import (
	"context"

	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/request"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/response"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error)
	Logout(ctx context.Context, req request.LogoutRequest) error
	RefreshToken(ctx context.Context, req request.RefreshTokenRequest) (*response.RefreshTokenResponse, error)
	GenerateToken(accountId int64, role string) (*string, error)
	GetTokenExpiry() int64
}
