output "ecr_repository_url" {
  description = "ECR Repository URL"
  value       = aws_ecr_repository.app_repo.repository_url
}

output "ecr_repository_arn" {
  description = "ECR Repository ARN"
  value       = aws_ecr_repository.app_repo.arn
}

output "eks_cluster_name" {
  description = "EKS Cluster Name"
  value       = data.aws_eks_cluster.existing_cluster.name
}

output "eks_cluster_endpoint" {
  description = "EKS Cluster Endpoint"
  value       = data.aws_eks_cluster.existing_cluster.endpoint
}

output "eks_cluster_version" {
  description = "EKS Cluster Version"
  value       = data.aws_eks_cluster.existing_cluster.version
}

output "eks_node_role_arn" {
  description = "EKS Node Group Role ARN"
  value       = data.aws_iam_role.existing_node_role.arn
}

output "github_actions_user_arn" {
  description = "GitHub Actions IAM User ARN"
  value       = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.github_actions_user}"
}

output "aws_auth_config_applied" {
  description = "AWS Auth ConfigMap applied status"
  value       = "GitHub Actions user ${var.github_actions_user} has been granted access to EKS cluster ${var.eks_cluster_name}"
}

output "dynamodb_table_name" {
  description = "DynamoDB Table Name for Terraform State Locking"
  value       = aws_dynamodb_table.terraform_state_lock.name
}

output "s3_bucket_name" {
  description = "S3 Bucket Name for Terraform State"
  value       = aws_s3_bucket.terraform_state.bucket
}

output "s3_bucket_arn" {
  description = "S3 Bucket ARN for Terraform State"
  value       = aws_s3_bucket.terraform_state.arn
}