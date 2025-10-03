# 간단한 Terraform EC2 배포

가장 기본적인 Terraform 구성으로 AWS EC2 인스턴스를 배포합니다.

## 포함된 리소스

- EC2 인스턴스 (t3.micro, Amazon Linux 2)
- 보안 그룹 (HTTP, SSH 허용)
- Apache 웹서버 자동 설치

## 사용 방법

### 1. AWS 자격 증명 설정
```bash
aws configure
```

### 2. Terraform 실행
```bash
cd terraform
terraform init
terraform plan
terraform apply
```

### 3. 웹사이트 확인
출력된 `website_url`로 접속해서 웹페이지 확인

### 4. 정리
```bash
terraform destroy
```

## 특징

- **간단함**: 최소한의 설정으로 동작
- **기본 VPC 사용**: 새로운 네트워크 생성 없음
- **프리티어**: t3.micro 인스턴스 사용
- **자동 설치**: Apache 웹서버 자동 구성