# S3 Backend 설정 완료

## 현재 설정

✅ **S3 버킷**: `terraform-state-0ss4kx0a`
✅ **DynamoDB 테이블**: `my-terraform-project-terraform-state-lock`
✅ **리전**: `us-east-2`

## Backend 설정 파일

`backend.tf` 파일이 이미 생성되어 있습니다:

```hcl
terraform {
  backend "s3" {
    bucket         = "terraform-state-0ss4kx0a"
    key            = "terraform.tfstate"
    region         = "us-east-2"
    dynamodb_table = "my-terraform-project-terraform-state-lock"
    encrypt        = true
  }
}
```

## 사용 방법

```bash
cd terraform
terraform init    # S3 backend 초기화
terraform plan    # 변경사항 확인
terraform apply   # 인프라 배포
```

## 특징

- **고정 버킷**: 매번 같은 S3 버킷 사용
- **State 공유**: 팀원들과 동일한 State 파일 공유
- **버전 관리**: S3 버전 관리로 State 히스토리 보존
- **동시 실행 방지**: DynamoDB Lock으로 충돌 방지