# 기존 DynamoDB 테이블을 Terraform State로 가져오기 위한 리소스 정의
# Import 후에는 이 파일을 삭제해도 됩니다.

resource "aws_dynamodb_table" "existing_lock_table" {
  name           = "my-terraform-project-terraform-state-lock"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "My Terraform Project State Lock"
    Environment = "infrastructure"
    Project     = "my-terraform-project"
    ManagedBy   = "terraform"
  }
}

# Import 명령어:
# terraform import aws_dynamodb_table.existing_lock_table my-terraform-project-terraform-state-lock