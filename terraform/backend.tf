terraform {
  backend "s3" {
    bucket         = "terraform-state-0ss4kx0a"
    key            = "terraform.tfstate"
    region         = "us-east-2"
    dynamodb_table = "terraform-state-lock-0ss4kx0a"
    encrypt        = true
  }
}