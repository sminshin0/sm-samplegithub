terraform {
  backend "s3" {
    bucket         = "terraform-state-0ss4kx0a"
    key            = "terraform.tfstate"
    region         = "us-east-2"
    dynamodb_table = "my-terraform-project-terraform-state-lock"
    encrypt        = true
  }
}