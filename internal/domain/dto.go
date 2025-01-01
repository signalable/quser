package domain

import "time"

// 회원가입 요청 DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

// 프로필 업데이트 요청 DTO
type UpdateProfileRequest struct {
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
}

// 사용자 응답 DTO
type UserResponse struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Name       string       `json:"name"`
	Status     string       `json:"status"`
	IsVerified bool         `json:"is_verified"`
	Profile    *UserProfile `json:"profile,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}
