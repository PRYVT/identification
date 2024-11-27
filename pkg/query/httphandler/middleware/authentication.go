package middleware

import (
	"net/http"

	"github.com/L4B0MB4/PRYVT/identification/pkg/query/authentication"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthMiddleware struct {
	tokenManager *authentication.TokenManager
}

func NewAuthMiddleware(tm *authentication.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tm}
}

func (am *AuthMiddleware) AuthenticateMiddleware(c *gin.Context) {
	// Retrieve the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Error().Msg("Authorization header is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Check if the token has the Bearer prefix
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		log.Error().Msg("Authorization header format must be Bearer {token}")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Extract the token
	tokenString := authHeader[len(prefix):]

	// Verify the token
	_, err := am.tokenManager.VerifyToken(tokenString)
	if err != nil {
		log.Error().Err(err).Msg("Token verification failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
