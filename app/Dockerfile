# 멀티 스테이지 빌드로 최적화
FROM golang:1.21-alpine AS builder

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사
COPY go.mod go.sum* ./

# 의존성 다운로드 (캐시 최적화)
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드 (정적 링크)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 실행 스테이지 (경량화)
FROM alpine:latest

# 보안을 위한 CA 인증서 추가
RUN apk --no-cache add ca-certificates

# 작업 디렉토리 설정
WORKDIR /root/

# 빌드된 바이너리 복사
COPY --from=builder /app/main .

# 포트 노출
EXPOSE 8080

# 헬스체크 추가
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 애플리케이션 실행
CMD ["./main"]