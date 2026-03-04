package usecase

import (
	"context"
	"fmt"
	"irfanard27/incore-api/internal/domain/entity"
	"irfanard27/incore-api/internal/domain/repository"
	"irfanard27/incore-api/internal/infra/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	Logout(ctx context.Context, token string) error
}

type authUsecase struct {
	userRepo   repository.UserRepository
	jwtService jwt.JWTService
}

func NewAuthUsecase(userRepo repository.UserRepository, jwtService jwt.JWTService) AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (string, *entity.User, error) {

	// Get user by email
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, fmt.Errorf("user not found")
	}

	// Compare provided password with stored bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT access token
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, user, nil
}

func (a *authUsecase) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to check existing user: %w", err)
	// }
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user with hashed password
	user.Password = string(hashedPassword)

	// Create user
	err = a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (a *authUsecase) Logout(ctx context.Context, token string) error {
	// In a stateless JWT implementation, logout typically means:
	// 1. Add token to blacklist (if using token blacklisting)
	// 2. Remove refresh token from storage (if storing refresh tokens)
	// 3. Client-side token deletion

	// For now, this is a placeholder implementation
	// In production, you might want to implement token blacklisting
	if token == "" {
		return fmt.Errorf("token is required")
	}

	// Validate the token first to ensure it's valid before "logging out"
	_, err := a.jwtService.ValidateAccessToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Token validation successful - logout logic would go here
	// For now, we'll just return success
	return nil
}
