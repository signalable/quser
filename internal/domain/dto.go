package domain

import "time"

// LoginRequest 로그인 요청 DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse 로그인 응답 DTO
type LoginResponse struct {
	AccessToken string        `json:"access_token"`
	User        *UserResponse `json:"user"`
}

// RegisterRequest 회원가입 요청 DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

// UpdateProfileRequest 프로필 업데이트 요청 DTO
type UpdateProfileRequest struct {
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
}

// UserResponse 사용자 응답 DTO
type UserResponse struct {
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Name       string       `json:"name"`
	Status     string       `json:"status"`
	IsVerified bool         `json:"is_verified"`
	Profile    *UserProfile `json:"profile,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}
