# Backend 설정 가이드

## 1단계: Backend 인프라 생성

```bash
cd terraform
terraform init
terraform apply
```

출력된 S3 버킷 이름과 DynamoDB 테이블 이름을 기록해두세요.

## 2단계: main.tf에 backend 설정 추가

출력된 정보로 main.tf의 terraform 블록을 다음과 같이 수정:

```hcl
terraform {
  backend "s3" {
    bucket         = "terraform-state-xxxxxxxx"  # 출력된 버킷 이름
    key            = "terraform.tfstate"
    region         = "us-east-2"
    dynamodb_table = "terraform-state-lock-xxxxxxxx"  # 출력된 테이블 이름
    encrypt        = true
  }
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
  required_version = ">= 1.0"
}
```

## 3단계: backend-setup.tf 파일 삭제

Backend 설정 완료 후 `backend-setup.tf` 파일을 삭제하세요.

## 4단계: State 마이그레이션

```bash
terraform init  # backend 설정 적용
```

"Do you want to copy existing state to the new backend?" → **yes** 입력