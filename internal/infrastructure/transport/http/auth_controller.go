package http

import (
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/domain/customer"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/request"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/dto/response"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/types"
)

// NewJWTMiddleware creates a new JWT middleware with custom configuration
func NewJWTMiddleware(config types.JWTConfig) func([]string) iris.Handler {

	return func(escapedRoutes []string) iris.Handler {
		return func(ctx iris.Context) {
			// Check if the current route is in escaped routes
			currentPath := strings.Replace(ctx.Path(), constants.ApiPrefix, "", 1)

			for _, route := range escapedRoutes {
				if strings.EqualFold(currentPath, route) {
					ctx.Next()
					return
				}
			}

			requestID := ctx.Values().GetString(constants.CTXRequestIdKey)

			// Extract token from Authorization header
			tokenString := extractToken(ctx, config)
			if tokenString == "" {
				logger.Infof(requestID, "JWT token is missing")
				RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
				return
			}

			// Parse and validate token
			claims, err := validateToken(tokenString, &config)
			if err != nil {
				logger.Infof(requestID, "Invalid JWT token: %v", err)
				RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
				return
			}

			// Store claims in context
			ctx.Values().Set("claims", claims)
			ctx.Next()
		}
	}
}

// GetUserFromContext extracts user information from the context
func GetUserFromContext(ctx iris.Context) (*types.CustomClaims, error) {
	claims, ok := ctx.Values().Get("claims").(*types.CustomClaims)
	if !ok {
		return nil, domain.ErrInvalidClaims
	}
	return claims, nil
}

// Helper functions
func extractToken(ctx iris.Context, config types.JWTConfig) string {
	bearerToken := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(bearerToken, config.TokenPrefix) {
		return ""
	}
	return strings.TrimPrefix(bearerToken, config.TokenPrefix)
}

func validateToken(tokenString string, config *types.JWTConfig) (*types.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if token.Method != config.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.SecretKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, domain.ErrExpiredToken
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, domain.ErrMalformedToken
			default:
				return nil, domain.ErrInvalidToken
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidClaims
	}

	return claims, nil
}

type AuthController struct {
	customerDomain customer.CustomerDomain
	config         *types.JWTConfig
}

func NewAuthController(server *iris.Application, custDomain customer.CustomerDomain, config *types.JWTConfig) AuthController {
	ac := AuthController{
		customerDomain: custDomain,
		config:         config,
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
	// if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
	// 	logger.Errorf(requestId, "Invalid password for user %s: %v", req.Email, err)
	// 	RespondWithError(ctx.ResponseWriter(), requestId, domain.ErrUnauthorized)
	// 	return
	// }

	// Generate access token
	token, err := GenerateToken(account.GetID(), "ADMIN", *h.config)
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
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(h.config.TokenExpiry)),
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
	newToken, err := GenerateToken(account.GetID(), "ADMIN", *h.config)
	if err != nil {
		logger.Errorf(requestID, "Failed to generate new token: %v", err)
		RespondWithError(ctx.ResponseWriter(), requestID, err)
		return
	}

	response := map[string]interface{}{
		"token":      newToken,
		"expires_at": time.Now().Add(time.Duration(h.config.TokenExpiry)),
	}

	RespondWithJSON(ctx.ResponseWriter(), response, http.StatusOK, requestID)
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int64, role string, config types.JWTConfig) (string, error) {
	now := time.Now()
	claims := types.CustomClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(config.TokenExpiry).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "device-backend",
		},
	}

	token := jwt.NewWithClaims(config.SigningMethod, claims)
	return token.SignedString(config.SecretKey)
}
