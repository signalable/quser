// quser/internal/repository/interfaces.go
package repository

import (
	"context"

	"github.com/signalable/quser/internal/domain"
)

// UserRepository 인터페이스 정의
type UserRepository interface {
	// 사용자 생성
	Create(ctx context.Context, user *domain.User) error
	// 이메일로 사용자 찾기
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	// ID로 사용자 찾기
	FindByID(ctx context.Context, id string) (*domain.User, error)
	// 이메일 존재 여부 확인
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// 사용자 정보 업데이트
	Update(ctx context.Context, user *domain.User) error
	// 프로필 업데이트
	UpdateProfile(ctx context.Context, userID string, profile *domain.UserProfile) error
	// 이메일 인증 상태 업데이트
	UpdateVerificationStatus(ctx context.Context, userID string, isVerified bool) error
	// 사용자 상태 업데이트
	UpdateStatus(ctx context.Context, userID string, status string) error
}
