package http

import (
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/domain/customer"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/request"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/response"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/service"
)

type AuthController struct {
	customerDomain customer.CustomerDomain
	authService    service.AuthService
}

func NewAuthController(server *iris.Application, custDomain customer.CustomerDomain) AuthController {
	ac := AuthController{
		customerDomain: custDomain,
	}
	server.Post(constants.ApiPrefix+"/login", ac.HandleLogin)
	server.Post(constants.ApiPrefix+"/logout", ac.HandleLogout)
	server.Post(constants.ApiPrefix+"/refresh", ac.HandleRefreshToken)
	return ac
}

func (h *AuthController) HandleLogin(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	var req request.LoginRequest

	// Validate request body
	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestId)
		return
	}

	// Get user from database
	account, err := h.customerDomain.RetrieveAccount(requestId, req.Email)
	if err != nil {
		logger.Errorf(requestId, "User not found: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestId, domain.ErrUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		logger.Errorf(requestId, "Invalid password for user %s: %v", req.Email, err)
		RespondWithError(ctx.ResponseWriter(), requestId, domain.ErrUnauthorized)
		return
	}

	// Generate access token
	token, err := h.authService.GenerateToken(account.GetID(), "ADMIN")
	if err != nil {
		logger.Errorf(requestId, "Failed to generate token: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		logger.Errorf(requestId, "Failed to generate refresh token: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	response := response.LoginResponse{
		Token:        *token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(h.authService.GetTokenExpiry())),
		User: response.UserInfo{
			ID:        account.GetID(),
			Email:     account.GetEmail(),
			Name:      account.GetName(),
			Role:      "ADMIN",
			CreatedAt: account.GetCreatedAt(),
		},
	}

	RespondWithJSON(ctx.ResponseWriter(), response, http.StatusCreated, requestId)
}

// Logout handles user logout by invalidating the refresh token
func (h *AuthController) HandleLogout(ctx iris.Context) {
	requestID := ctx.Values().GetString(constants.CTXRequestIdKey)

	// Get refresh token from request
	refreshToken := ctx.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		logger.Infof(requestID, "No refresh token provided for logout")
		RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), map[string]string{
		"message": "Successfully logged out",
	}, http.StatusOK, requestID)
}

// HandleRefreshToken handles token refresh requests
func (h *AuthController) HandleRefreshToken(ctx iris.Context) {
	requestID := ctx.Values().GetString(constants.CTXRequestIdKey)

	refreshToken := ctx.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		logger.Infof(requestID, "No refresh token provided")
		RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
		return
	}
	var req request.LoginRequest

	// Validate request body
	if err := GetRequest(ctx.Request(), &req); err != nil {
		RespondWithMappingError(ctx.ResponseWriter(), err.Error(), requestID)
		return
	}

	account, err := h.customerDomain.RetrieveAccount(requestID, req.Email)
	if err != nil {
		logger.Errorf(requestID, "User not found: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
		return
	}

	// Generate new access token
	newToken, err := h.authService.GenerateToken(account.GetID(), "ADMIN")
	if err != nil {
		logger.Errorf(requestID, "Failed to generate new token: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestID, err)
		return
	}

	response := map[string]interface{}{
		"token":      newToken,
		"expires_at": time.Now().Add(time.Duration(h.authService.GetTokenExpiry())),
	}

	RespondWithJSON(ctx.ResponseWriter(), response, http.StatusOK, requestID)
}
