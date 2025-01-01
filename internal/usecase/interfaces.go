package usecase

import (
	"context"

	"github.com/signalable/quser/internal/domain"
)

// UserUseCase 인터페이스 정의
type UserUseCase interface {
	// 회원가입
	Register(ctx context.Context, req *domain.RegisterRequest) error
	// 로그인
	Login(ctx context.Context, email, password string) (*domain.LoginResponse, error)
	// 프로필 조회
	GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error)
	// 프로필 업데이트
	UpdateProfile(ctx context.Context, userID string, req *domain.UpdateProfileRequest) error
	// 이메일 인증
	VerifyEmail(ctx context.Context, userID string, token string) error
	// 사용자 상태 조회
	GetUserStatus(ctx context.Context, userID string) (string, error)
	// 이메일로 사용자 찾기
	FindByEmail(ctx context.Context, email string) (*domain.UserResponse, error)
	// 로그아웃웃
	Logout(ctx context.Context, token string) error
}
