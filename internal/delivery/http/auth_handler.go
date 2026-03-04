package http

import (
	"fmt"
	"irfanard27/incore-api/internal/domain/dto"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/usecase"
	httpresponse "irfanard27/incore-api/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	token, user, err := h.authUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusUnauthorized)
		return
	}

	httpresponse.SuccessResponse(c, "Login successfully", dto.LoginResponse{
		AccessToken: token,
		User:        dto.ToUserDto(user),
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in usecase
	}

	createdUser, err := h.authUsecase.Register(c.Request.Context(), user)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	httpresponse.SuccessResponse(c, "Register successfully", dto.RegisterResponse{
		User: dto.ToUserDto(createdUser),
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		httpresponse.ErrorResponse(c, fmt.Errorf("authorization header is required"), http.StatusBadRequest)
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.authUsecase.Logout(c.Request.Context(), token)
	if err != nil {
		httpresponse.ErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	httpresponse.SuccessResponse(c, "Logged out successfully", nil)
}
