// quser/internal/usecase/user_usecase.go
package usecase

import (
	"context"
	"time"

	"github.com/signalable/quser/internal/client"
	"github.com/signalable/quser/internal/domain"
	"github.com/signalable/quser/internal/repository"
)

type userUseCase struct {
	userRepo   repository.UserRepository
	authClient *client.AuthClient
}

// NewUserUseCase User 유스케이스 생성자
func NewUserUseCase(
	userRepo repository.UserRepository,
	authClient *client.AuthClient,
) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		authClient: authClient,
	}
}

// Register 회원가입 구현
func (uc *userUseCase) Register(ctx context.Context, req *domain.RegisterRequest) error {
	// 이메일 중복 체크
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrEmailAlreadyExists
	}

	// 새 사용자 생성
	user := &domain.User{
		Email:  req.Email,
		Name:   req.Name,
		Status: domain.UserStatusPending,
		Profile: &domain.UserProfile{
			LastLogin: time.Now(),
		},
	}

	// 사용자 저장
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

// Login 로그인 구현
func (uc *userUseCase) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	// 사용자 조회
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Auth Service에 토큰 생성 요청
	authResp, err := uc.authClient.CreateToken(ctx, user.ID.Hex())
	if err != nil {
		return nil, err
	}

	// 로그인 응답 생성
	return &domain.LoginResponse{
		AccessToken: authResp.AccessToken,
		User: &domain.UserResponse{
			ID:         user.ID.Hex(),
			Email:      user.Email,
			Name:       user.Name,
			Status:     user.Status,
			IsVerified: user.IsVerified,
			Profile:    user.Profile,
			CreatedAt:  user.CreatedAt,
		},
	}, nil
}

// GetProfile 프로필 조회 구현
func (uc *userUseCase) GetProfile(ctx context.Context, userID string) (*domain.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:         user.ID.Hex(),
		Email:      user.Email,
		Name:       user.Name,
		Status:     user.Status,
		IsVerified: user.IsVerified,
		Profile:    user.Profile,
		CreatedAt:  user.CreatedAt,
	}, nil
}

// UpdateProfile 프로필 업데이트 구현
func (uc *userUseCase) UpdateProfile(ctx context.Context, userID string, req *domain.UpdateProfileRequest) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 프로필 정보 업데이트
	if user.Profile == nil {
		user.Profile = &domain.UserProfile{}
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.PhoneNumber != "" {
		user.Profile.PhoneNumber = req.PhoneNumber
	}
	if req.Bio != "" {
		user.Profile.Bio = req.Bio
	}
	if req.Avatar != "" {
		user.Profile.Avatar = req.Avatar
	}

	return uc.userRepo.UpdateProfile(ctx, userID, user.Profile)
}

// VerifyEmail 이메일 인증 구현
func (uc *userUseCase) VerifyEmail(ctx context.Context, userID string, token string) error {
	// 토큰 유효성 검증
	if err := uc.authClient.ValidateToken(ctx, token); err != nil {
		return err
	}

	// 이메일 인증 상태 업데이트
	if err := uc.userRepo.UpdateVerificationStatus(ctx, userID, true); err != nil {
		return err
	}

	return nil
}

// GetUserStatus 사용자 상태 조회 구현
func (uc *userUseCase) GetUserStatus(ctx context.Context, userID string) (string, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	return user.Status, nil
}

// FindByEmail 이메일로 사용자 찾기 구현
func (uc *userUseCase) FindByEmail(ctx context.Context, email string) (*domain.UserResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:         user.ID.Hex(),
		Email:      user.Email,
		Name:       user.Name,
		Status:     user.Status,
		IsVerified: user.IsVerified,
		Profile:    user.Profile,
		CreatedAt:  user.CreatedAt,
	}, nil
}

// Logout 로그아웃 구현
func (uc *userUseCase) Logout(ctx context.Context, token string) error {
	return uc.authClient.RevokeToken(ctx, token)
}
