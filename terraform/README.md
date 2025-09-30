# Terraform EC2 배포

이 Terraform 구성은 AWS에 EC2 인스턴스를 배포합니다.

## 포함된 리소스

- VPC (Virtual Private Cloud)
- Public Subnet
- Internet Gateway
- Route Table
- Security Group (SSH, HTTP 허용)
- EC2 Instance (Amazon Linux 2 + Apache)

## 사용 방법

### 1. 사전 준비
```bash
# AWS CLI 설정 (AWS 자격 증명 필요)
aws configure

# Terraform 설치 확인
terraform --version
```

### 2. 설정 파일 준비
```bash
# terraform.tfvars 파일 생성
cp terraform.tfvars.example terraform.tfvars

# terraform.tfvars 파일 편집 (키 페어 이름 등 설정)
```

### 3. 배포 실행
```bash
# Terraform 초기화
terraform init

# 실행 계획 확인
terraform plan

# 리소스 배포
terraform apply
```

### 4. 정리
```bash
# 리소스 삭제
terraform destroy
```

## 주요 변수

- `aws_region`: AWS 리전 (기본값: ap-northeast-2)
- `project_name`: 프로젝트 이름 (기본값: terraform-ec2)
- `instance_type`: EC2 인스턴스 타입 (기본값: t3.micro)
- `key_pair_name`: SSH 접속용 키 페어 이름

## 출력값

- `instance_public_ip`: EC2 인스턴스 공개 IP
- `instance_public_dns`: EC2 인스턴스 공개 DNS
- `instance_id`: EC2 인스턴스 ID

## 접속 방법

```bash
# SSH 접속 (키 페어 필요)
ssh -i your-key.pem ec2-user@<public_ip>

# 웹 브라우저에서 확인
http://<public_ip>
```