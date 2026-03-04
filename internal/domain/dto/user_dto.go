package dto

import (
	"time"

	"irfanard27/incore-api/internal/domain/entity"
)

type UserDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserDto(u *entity.User) *UserDTO {
	createdAt, _ := time.Parse(time.RFC3339, u.CreatedAt)
	return &UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: createdAt,
	}
}
