// quser/cmd/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/signalable/quser/internal/client"
	"github.com/signalable/quser/internal/config"
	"github.com/signalable/quser/internal/delivery/http/handler"
	"github.com/signalable/quser/internal/delivery/http/middleware"
	"github.com/signalable/quser/internal/delivery/http/routes"
	"github.com/signalable/quser/internal/repository/mongodb"
	"github.com/signalable/quser/internal/usecase"
)

func main() {
	// 설정 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("설정을 로드할 수 없습니다: %v", err)
	}

	// MongoDB 연결
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatalf("MongoDB 연결 실패: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// MongoDB 연결 테스트
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB 연결 테스트 실패: %v", err)
	}

	// Auth 클라이언트 초기화
	authClient := client.NewAuthClient(
		cfg.AuthService.URL,
		cfg.AuthService.Timeout,
	)

	// 레포지토리 초기화
	userRepo := mongodb.NewUserRepository(mongoClient.Database(cfg.MongoDB.Database))

	// 유스케이스 초기화
	userUseCase := usecase.NewUserUseCase(userRepo, authClient)

	// 핸들러 및 미들웨어 초기화
	userHandler := handler.NewUserHandler(userUseCase)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// 라우터 설정
	router := mux.NewRouter()
	routes.SetupUserRoutes(router, userHandler, authMiddleware)

	// CORS 미들웨어 설정
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// 서버 시작
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("User Service 시작: %s", serverAddr)

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
