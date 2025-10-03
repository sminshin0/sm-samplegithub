# Terraform 변수 정의

variable "aws_region" {
  description = "AWS 리전"
  type        = string
  default     = "us-east-2"
}

variable "project_name" {
  description = "프로젝트 이름 (리소스 태그에 사용)"
  type        = string
  default     = "simple-terraform"
}

variable "instance_type" {
  description = "EC2 인스턴스 타입"
  type        = string
  default     = "t3.micro"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}