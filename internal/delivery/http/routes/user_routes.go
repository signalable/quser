package routes

import (
	"github.com/gorilla/mux"
	"github.com/signalable/quser/internal/delivery/http/handler"
	"github.com/signalable/quser/internal/delivery/http/middleware"
)

// SetupUserRoutes 라우터 설정
func SetupUserRoutes(
	router *mux.Router,
	userHandler *handler.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// 공개 라우트
	router.HandleFunc("/api/users/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/users/login", userHandler.Login).Methods("POST")

	// 인증이 필요한 라우트
	router.HandleFunc("/api/users/logout", authMiddleware.Authenticate(userHandler.Logout)).Methods("POST")
	router.HandleFunc("/api/users/{id}/profile", authMiddleware.Authenticate(userHandler.GetProfile)).Methods("GET")
	router.HandleFunc("/api/users/{id}/profile", authMiddleware.Authenticate(userHandler.UpdateProfile)).Methods("PUT")
}
