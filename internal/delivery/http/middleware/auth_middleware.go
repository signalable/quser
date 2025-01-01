package middleware

import (
	"net/http"
	"strings"

	"github.com/signalable/quser/internal/client"
)

type AuthMiddleware struct {
	authClient *client.AuthClient
}

// NewAuthMiddleware Auth 미들웨어 생성자
func NewAuthMiddleware(authClient *client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Authenticate 인증 미들웨어
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "인증이 필요합니다", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "잘못된 인증 형식입니다", http.StatusUnauthorized)
			return
		}

		if err := m.authClient.ValidateToken(r.Context(), tokenParts[1]); err != nil {
			http.Error(w, "유효하지 않은 토큰입니다", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
