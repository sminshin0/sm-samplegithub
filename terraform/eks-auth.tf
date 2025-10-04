# 현재 AWS 계정 정보 가져오기
data "aws_caller_identity" "current" {}

# EKS 클러스터 인증 정보
data "aws_eks_cluster_auth" "main" {
  name = aws_eks_cluster.main.name
}

# Kubernetes provider 설정 (재시도 로직 포함)
provider "kubernetes" {
  host                   = aws_eks_cluster.main.endpoint
  cluster_ca_certificate = base64decode(aws_eks_cluster.main.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.main.token
  
  # 연결 재시도 설정
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args = [
      "eks",
      "get-token",
      "--cluster-name",
      aws_eks_cluster.main.name,
      "--region",
      var.aws_region
    ]
  }
}

# EKS 클러스터 준비 대기
resource "time_sleep" "wait_for_cluster" {
  depends_on = [
    aws_eks_cluster.main,
    aws_eks_node_group.main
  ]
  
  create_duration = "30s"
}

# EKS aws-auth 설정 (eksctl 사용 - 더 안전함)
resource "null_resource" "aws_auth" {
  count = var.create_aws_auth ? 1 : 0
  
  provisioner "local-exec" {
    command = <<-EOT
      # GitHub Actions 사용자 권한 추가
      eksctl create iamidentitymapping \
        --cluster ${aws_eks_cluster.main.name} \
        --region ${var.aws_region} \
        --arn arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.github_actions_user} \
        --group system:masters \
        --username github-actions \
        --no-duplicate-arns || echo "IAM identity mapping may already exist"
    EOT
  }

  depends_on = [
    time_sleep.wait_for_cluster
  ]

  triggers = {
    cluster_name = aws_eks_cluster.main.name
    user_arn     = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.github_actions_user}"
  }
}