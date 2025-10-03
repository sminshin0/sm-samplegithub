# 기본 Terraform 설정
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.0"

  backend "s3" {
    bucket         = "terraform-state-0ss4kx0a"
    key            = "terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "my-terraform-project-terraform-state-lock"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region
}

# 기존 EKS 클러스터 참조 (수동으로 생성됨)
data "aws_eks_cluster" "existing_cluster" {
  name = "sm-eks"
}

data "aws_eks_cluster_auth" "existing_cluster" {
  name = "sm-eks"
}

# 기존 EKS 노드 역할 참조 (수동으로 생성됨)
data "aws_iam_role" "existing_node_role" {
  name = "github-actions-terraform-eks-node-role"
}

# S3 Bucket for Terraform State
resource "aws_s3_bucket" "terraform_state" {
  bucket = "terraform-state-0ss4kx0a"

  tags = {
    Name        = "${var.project_name}-terraform-state"
    Environment = var.environment
    ManagedBy   = "terraform"
    Purpose     = "terraform-state-storage"
  }
}

# S3 Bucket Versioning
resource "aws_s3_bucket_versioning" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 Bucket Server Side Encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# S3 Bucket Public Access Block
resource "aws_s3_bucket_public_access_block" "terraform_state" {
  bucket = aws_s3_bucket.terraform_state.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# DynamoDB Table for Terraform State Locking
resource "aws_dynamodb_table" "terraform_state_lock" {
  name           = "my-terraform-project-terraform-state-lock"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "${var.project_name}-terraform-lock"
    Environment = var.environment
    ManagedBy   = "terraform"
    Purpose     = "terraform-state-locking"
  }
}

# ECR Repository
resource "aws_ecr_repository" "app_repo" {
  name                 = "${var.project_name}-app"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = "${var.project_name}-ecr"
    Environment = var.environment
    ManagedBy   = "terraform"
  