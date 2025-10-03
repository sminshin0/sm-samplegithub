# DynamoDB 테이블 생성용 임시 파일
# 이 파일로 테이블을 생성한 후 삭제할 예정

resource "aws_dynamodb_table" "terraform_state_lock" {
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
  }
}