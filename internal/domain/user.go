package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 상태 상수
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusPending  = "pending"
)

// User 도메인 모델
type User struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"-" bson:"password"` // JSON 직렬화에서 제외
	Name       string             `json:"name" bson:"name"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	IsVerified bool               `json:"is_verified" bson:"is_verified"`
	Profile    *UserProfile       `json:"profile,omitempty" bson:"profile,omitempty"`
}

// UserProfile 도메인 모델
type UserProfile struct {
	PhoneNumber string    `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Avatar      string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Bio         string    `json:"bio,omitempty" bson:"bio,omitempty"`
	LastLogin   time.Time `json:"last_login,omitempty" bson:"last_login,omitempty"`
}
