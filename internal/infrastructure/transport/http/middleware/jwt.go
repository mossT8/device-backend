package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/application/logger"
	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
)

// JWTConfig holds the configuration for JWT middleware
type JWTConfig struct {
	SecretKey     []byte
	TokenExpiry   time.Duration
	SigningMethod jwt.SigningMethod
	TokenPrefix   string
}

// CustomClaims extends jwt.StandardClaims to include user-specific claims
type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.StandardClaims
}

var (
	defaultConfig = JWTConfig{
		SecretKey:     []byte("super_duper_secret_key"), // Should be loaded from environment variables
		TokenExpiry:   72 * time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
		TokenPrefix:   "Bearer ",
	}
)

// NewJWTMiddleware creates a new JWT middleware with custom configuration
func NewJWTMiddleware(config *JWTConfig) func([]string) iris.Handler {
	if config == nil {
		config = &defaultConfig
	}

	return func(escapedRoutes []string) iris.Handler {
		return func(ctx iris.Context) {
			// Check if the current route is in escaped routes
			currentPath := ctx.Path()
			for _, route := range escapedRoutes {
				if strings.HasPrefix(currentPath, route) {
					ctx.Next()
					return
				}
			}

			requestID := ctx.Values().GetString(constants.CTXRequestIdKey)

			// Extract token from Authorization header
			tokenString := extractToken(ctx)
			if tokenString == "" {
				logger.Infof(requestID, "JWT token is missing")
				http.RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
				return
			}

			// Parse and validate token
			claims, err := validateToken(tokenString, config)
			if err != nil {
				logger.Infof(requestID, "Invalid JWT token: %v", err)
				http.RespondWithError(ctx.ResponseWriter(), requestID, domain.ErrUnauthorized)
				return
			}

			// Store claims in context
			ctx.Values().Set("claims", claims)
			ctx.Next()
		}
	}
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int64, role string, config *JWTConfig) (string, error) {
	if config == nil {
		config = &defaultConfig
	}

	now := time.Now()
	claims := CustomClaims{
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

// GetUserFromContext extracts user information from the context
func GetUserFromContext(ctx iris.Context) (*CustomClaims, error) {
	claims, ok := ctx.Values().Get("claims").(*CustomClaims)
	if !ok {
		return nil, domain.ErrInvalidClaims
	}
	return claims, nil
}

// Helper functions

func extractToken(ctx iris.Context) string {
	bearerToken := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(bearerToken, defaultConfig.TokenPrefix) {
		return ""
	}
	return strings.TrimPrefix(bearerToken, defaultConfig.TokenPrefix)
}

func validateToken(tokenString string, config *JWTConfig) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidClaims
	}

	return claims, nil
}
