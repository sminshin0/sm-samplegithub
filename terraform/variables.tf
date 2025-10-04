variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "github-actions-terraform"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "github_actions_user" {
  description = "GitHub Actions IAM user name"
  type        = string
  default     = "sm-user"
}

variable "eks_cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "sm-eks"
}