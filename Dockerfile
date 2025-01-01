# 빌드 스테이지
FROM golang:1.22-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일 복사
COPY go.mod go.sum ./

# 의존성 다운로드
RUN go mod download

# 소스 코드 복사
COPY . .

# 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 실행 스테이지
FROM alpine:latest

# 필요한 CA certificates 설치
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 빌드 스테이지에서 바이너리 복사
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# 포트 설정 (Auth Service는 8080, User Service는 8081 사용)
EXPOSE 8081

# 실행
CMD ["./main"]