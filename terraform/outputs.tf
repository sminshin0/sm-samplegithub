# 중요한 정보들을 출력
output "instance_id" {
  description = "EC2 인스턴스 ID"
  value       = aws_instance.web.id
}

output "public_ip" {
  description = "EC2 인스턴스 공개 IP"
  value       = aws_instance.web.public_ip
}

output "public_dns" {
  description = "EC2 인스턴스 공개 DNS"
  value       = aws_instance.web.public_dns
}

output "website_url" {
  description = "웹사이트 URL"
  value       = "http://${aws_instance.web.public_ip}"
}

output "ecr_repository_url" {
  description = "ECR Repository URL"
  value       = aws_ecr_repository.app_repo.repository_url
}

output "eks_cluster_name" {
  description = "EKS Cluster Name"
  value       = aws_eks_cluster.main.name
}

output "eks_cluster_endpoint" {
  description = "EKS Cluster Endpoint"
  value       = aws_eks_cluster.main.endpoint
}

output "eks_cluster_arn" {
  description = "EKS Cluster ARN"
  value       = aws_eks_cluster.main.arn
}