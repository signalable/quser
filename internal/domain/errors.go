// quser/internal/domain/errors.go
package domain

import "errors"

var (
	// 사용자 관련 에러
	ErrUserNotFound       = errors.New("사용자를 찾을 수 없습니다")
	ErrEmailAlreadyExists = errors.New("이미 존재하는 이메일입니다")
	ErrInvalidCredentials = errors.New("잘못된 인증 정보입니다")

	// 프로필 관련 에러
	ErrInvalidProfileData = errors.New("잘못된 프로필 데이터입니다")

	// 검증 관련 에러
	ErrEmailVerification = errors.New("이메일 검증에 실패했습니다")
)
