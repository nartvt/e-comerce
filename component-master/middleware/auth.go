package middleware

import (
	cfg "component-master/config"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Bearer        = "Bearer "
	Authorization = "Authorization"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthHandler struct {
	config       cfg.TokenConfig
	AllowedPaths []string // Paths that don't require authentication
}

func NewAuthHandler(config cfg.TokenConfig) *AuthHandler {
	return &AuthHandler{
		config: config,
	}
}

func (h *AuthHandler) GenerateTokenPair(userID uint, email string) (*TokenPair, error) {
	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.config.AccessTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	accessTokenString, err := accessToken.SignedString([]byte(h.config.AccessTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.config.RefreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(h.config.RefreshTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// Middleware to verify access token
func (h *AuthHandler) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if path is allowed without auth
		path := c.Path()
		for _, allowedPath := range h.AllowedPaths {
			if strings.HasPrefix(path, allowedPath) {
				return c.Next()
			}
		}

		// Get token from Authorization header
		authHeader := c.Get(Authorization)
		if len(authHeader) <= 7 || authHeader[:7] != Bearer {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header",
			})
		}

		tokenString := authHeader[7:]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.config.AccessTokenSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid access token",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid access token",
			})
		}

		// Add claims to context for later use
		c.Locals("user", claims)

		return c.Next()
	}
}
