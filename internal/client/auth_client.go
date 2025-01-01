package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/signalable/quser/internal/domain"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TokenValidationResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id"`
}

// NewAuthClient Auth 클라이언트 생성자
func NewAuthClient(baseURL string, timeout time.Duration) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// CreateToken 토큰 생성 요청
func (c *AuthClient) CreateToken(ctx context.Context, userID string) (*AuthResponse, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/auth/token", c.baseURL),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("요청 생성 실패: %w", err)
	}

	// User ID를 헤더에 추가
	req.Header.Set("X-User-ID", userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("토큰 생성 실패: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	return &authResp, nil
}

// ValidateToken 토큰 검증
func (c *AuthClient) ValidateToken(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/api/auth/token/validate", c.baseURL),
		nil,
	)
	if err != nil {
		return fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("토큰 검증 실패: %d", resp.StatusCode)
	}

	var validationResp TokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
		return fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if !validationResp.Valid {
		return fmt.Errorf("유효하지 않은 토큰")
	}

	return nil
}

// RevokeToken 토큰 폐기 요청
func (c *AuthClient) RevokeToken(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/auth/token/revoke", c.baseURL),
		nil,
	)
	if err != nil {
		return fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return domain.ErrInvalidToken
		}
		return domain.ErrLogoutFailed
	}

	return nil
}
