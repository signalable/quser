package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
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

// ValidateToken 토큰 검증
func (c *AuthClient) ValidateToken(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/api/auth/validate", c.baseURL),
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

	return nil
}

// RequestPasswordHash 비밀번호 해시 요청
func (c *AuthClient) RequestPasswordHash(ctx context.Context, password string) (string, error) {
	type hashRequest struct {
		Password string `json:"password"`
	}

	payload, err := json.Marshal(hashRequest{Password: password})
	if err != nil {
		return "", fmt.Errorf("요청 데이터 직렬화 실패: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/api/auth/hash", c.baseURL),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return "", fmt.Errorf("요청 생성 실패: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("해시 요청 실패: %d", resp.StatusCode)
	}

	var result struct {
		Hash string `json:"hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("응답 파싱 실패: %w", err)
	}

	return result.Hash, nil
}
