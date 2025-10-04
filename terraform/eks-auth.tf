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

# EKS aws-auth ConfigMap 관리 (조건부)
resource "kubernetes_config_map" "aws_auth" {
  count = var.create_aws_auth ? 1 : 0
  
  metadata {
    name      = "aws-auth"
    namespace = "kube-system"
  }

  data = {
    mapUsers = yamlencode([
      {
        userarn  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:user/${var.github_actions_user}"
        username = "github-actions"
        groups   = ["system:masters"]
      }
    ])
    
    # 노드 그룹 역할 (자동 생성됨)
    mapRoles = yamlencode([
      {
        rolearn  = aws_iam_role.eks_node_role.arn
        username = "system:node:{{EC2PrivateDNSName}}"
        groups   = [
          "system:bootstrappers",
          "system:nodes"
        ]
      }
    ])
  }

  depends_on = [
    time_sleep.wait_for_cluster
  ]
}