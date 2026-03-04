package middleware

import (
	"fmt"
	"irfanard27/incore-api/internal/infra/jwt"
	httpresponse "irfanard27/incore-api/pkg/http_response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	jwtService jwt.JWTService
}

func NewJWTMiddleware(jwtService jwt.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

func (m *JWTMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpresponse.ErrorResponse(c, fmt.Errorf("authorization header required"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			httpresponse.ErrorResponse(c, fmt.Errorf("invalid authorization header format"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			httpresponse.ErrorResponse(c, fmt.Errorf("invalid or expired token"), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
