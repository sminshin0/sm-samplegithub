# 현재 AWS 계정 정보 가져오기
data "aws_caller_identity" "current" {}

# EKS aws-auth ConfigMap 관리
resource "kubernetes_config_map" "aws_auth" {
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
    
    # 기존 노드 그룹 역할도 유지
    mapRoles = yamlencode([
      {
        rolearn  = data.aws_iam_role.existing_node_role.arn
        username = "system:node:{{EC2PrivateDNSName}}"
        groups   = [
          "system:bootstrappers",
          "system:nodes"
        ]
      }
    ])
  }

  depends_on = [
    data.aws_eks_cluster.existing_cluster
  ]
}

# Kubernetes provider 설정
provider "kubernetes" {
  host                   = data.aws_eks_cluster.existing_cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.existing_cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.existing_cluster.token
}