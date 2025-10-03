# DynamoDB 테이블 생성용 임시 파일
# 이 파일로 테이블을 생성한 후 삭제할 예정

resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "terraform-state-lock-0ss4kx0a"
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