# S3 버킷과 DynamoDB 테이블을 생성하는 별도 파일
# 이 파일로 먼저 backend 인프라를 만들고, 나중에 삭제할 예정

# 랜덤 문자열로 고유한 버킷 이름 생성
resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Terraform State 저장용 S3 버킷
resource "aws_s3_bucket" "terraform_state" {
  bucket = "terraform-state-${random_string.bucket_suffix.result}"

  tags = {
    Name        = "Terraform State Bucket"
    Environment = "infrastructure"
  }
}

# S3 버킷 버전 관리 활성화
resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 버킷 암호화 설정
resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# S3 버킷 퍼블릭 액세스 차단
resource "aws_s3_bucket_public_access_block" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# DynamoDB 테이블 (State Lock용)
resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "terraform-state-lock-${random_string.bucket_suffix.result}"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "Terraform State Lock Table"
    Environment = "infrastructure"
  }
}

# 생성된 리소스 정보 출력
output "s3_bucket_name" {
  description = "Terraform State S3 버킷 이름"
  value       = aws_s3_bucket.terraform_state.bucket
}

output "dynamodb_table_name" {
  description = "Terraform State Lock DynamoDB 테이블 이름"
  value       = aws_dynamodb_table.terraform_state_lock.name
}