package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/signalable/quser/internal/domain"
	"github.com/signalable/quser/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler User 핸들러 생성자
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Register 회원가입 핸들러
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.Register(r.Context(), &req); err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "회원가입이 완료되었습니다",
	})
}

// Login 로그인 핸들러
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	resp, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetProfile 프로필 조회 핸들러
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	profile, err := h.userUseCase.GetProfile(r.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile 프로필 업데이트 핸들러
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req domain.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청 형식입니다", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.UpdateProfile(r.Context(), userID, &req); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case domain.ErrInvalidProfileData:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "프로필이 업데이트되었습니다",
	})
}

// VerifyEmail 이메일 인증 핸들러
func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "인증 토큰이 필요합니다", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.VerifyEmail(r.Context(), userID, token); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case domain.ErrEmailVerification:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "내부 서버 오류", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "이메일이 인증되었습니다",
	})

}

// Logout 로그아웃 핸들러
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "토큰이 없습니다", http.StatusUnauthorized)
		return
	}

	// "Bearer " 접두사 제거
	tokenString := token[7:]

	if err := h.userUseCase.Logout(r.Context(), tokenString); err != nil {
		http.Error(w, "로그아웃 처리 중 오류가 발생했습니다", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "로그아웃되었습니다",
	})
}
