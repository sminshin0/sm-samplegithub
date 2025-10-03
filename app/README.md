# Go Web Application

간단한 Go 웹 애플리케이션으로 "Hello World"를 출력합니다.

## 🚀 기능

- **웹 서버**: HTTP 서버로 HTML 페이지 제공
- **헬스체크**: `/health` 엔드포인트로 상태 확인
- **Docker 지원**: 멀티 스테이지 빌드로 최적화
- **자동 배포**: GitHub Actions로 ECR에 자동 푸시

## 📁 구조

```
app/
├── main.go          # 메인 애플리케이션
├── go.mod           # Go 모듈 정의
├── Dockerfile       # Docker 이미지 빌드
└── README.md        # 이 파일
```

## 🛠️ 로컬 실행

### Go로 직접 실행
```bash
cd app
go run main.go
```

### Docker로 실행
```bash
cd app
docker build -t go-hello-world .
docker run -p 8080:8080 go-hello-world
```

## 🌐 엔드포인트

- **메인 페이지**: http://localhost:8080/
- **헬스체크**: http://localhost:8080/health

## 🔄 CI/CD 워크플로우

1. **코드 변경**: `app/` 폴더 내 파일 수정
2. **자동 테스트**: Go 테스트 및 빌드 검증
3. **Docker 빌드**: 멀티 스테이지 빌드로 이미지 생성
4. **ECR 푸시**: AWS ECR에 이미지 자동 업로드
5. **태그 관리**: 버전별 태그 및 latest 태그 관리

## 📦 배포된 이미지 사용

```bash
# ECR에서 이미지 pull (AWS CLI 설정 필요)
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <ECR_URI>
docker pull <ECR_URI>:latest
docker run -p 8080:8080 <ECR_URI>:latest
```

## 🔧 환경 변수

- `PORT`: 서버 포트 (기본값: 8080)