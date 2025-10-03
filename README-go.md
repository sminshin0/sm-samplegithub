# Go Hello World 프로그램

## 파일 구조
```
├── main.go      # 기본 콘솔 Hello World
├── server.go    # 웹서버 Hello World  
├── go.mod       # Go 모듈 파일
└── README-go.md # 이 파일
```

## 실행 방법

### 1. 기본 Hello World (콘솔)
```bash
go run main.go
```

### 2. 웹서버 Hello World
```bash
go run server.go
```

그 후 브라우저에서 접속:
- http://localhost:8080/
- http://localhost:8080/health

### 3. 빌드해서 실행
```bash
# 콘솔 버전 빌드
go build -o hello main.go
./hello

# 웹서버 버전 빌드
go build -o server server.go
./server
```

## 기능

### main.go
- 간단한 "Hello, World!" 출력
- Go 언어의 가장 기본적인 프로그램

### server.go
- HTTP 웹서버 실행
- `/` 엔드포인트: Hello World + 현재 시간 + 요청자 IP
- `/health` 엔드포인트: 헬스체크용
- 포트 8080에서 실행

## Go 설치 확인
```bash
go version
```

Go가 설치되어 있지 않다면:
- macOS: `brew install go`
- Ubuntu: `sudo apt install golang-go`
- Windows: https://golang.org/dl/