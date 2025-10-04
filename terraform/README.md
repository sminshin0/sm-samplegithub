# Terraform Infrastructure

이 디렉토리는 AWS 인프라를 Terraform으로 관리합니다.

## 📁 구조

```
terraform/
├── main.tf          # 메인 Terraform 설정
├── variables.tf     # 변수 정의
├── outputs.tf       # 출력값 정의
├── eks-auth.tf      # EKS 인증 설정
└── README.md        # 이 파일
```

## 🚀 배포 방법

### Terraform 완전 자동화 배포

```bash
# 1. Terraform 디렉토리로 이동
cd terraform

# 2. 초기화 (최초 1회)
terraform init

# 3. 전체 인프라 배포 (약 15-20분 소요)
terraform apply

# 4. kubectl 설정 (배포 완료 후)
aws eks update-kubeconfig --region us-east-1 --name sm-eks
kubectl get nodes
```

**주의사항:**
- AWS CLI가 올바르게 설정되어 있어야 합니다
- `eksctl` 명령어가 설치되어 있어야 합니다 (aws-auth 설정용)
- Terraform이 EKS 클러스터를 생성하므로 권한 문제가 없습니다

**eksctl 설치:**
```bash
# macOS
brew install weaveworks/tap/eksctl

# 또는 다른 방법: https://eksctl.io/installation/
```

### 배포 순서
1. **S3 + DynamoDB**: Terraform state 관리
2. **ECR 리포지토리**: Docker 이미지 저장소
3. **IAM 역할들**: EKS 클러스터 및 노드 권한
4. **EKS 클러스터**: Kubernetes 마스터 노드 (~10분)
5. **EKS 노드 그룹**: 워커 노드들 (~5분)
6. **EKS 애드온**: VPC CNI, CoreDNS, kube-proxy
7. **aws-auth ConfigMap**: GitHub Actions 접근 권한



## 📦 생성되는 리소스

### **AWS 리소스**
- **S3 버킷**: Terraform state 저장
- **DynamoDB 테이블**: State 잠금 관리
- **ECR 리포지토리**: Docker 이미지 저장
- **EKS 클러스터**: Kubernetes 클러스터
- **EKS 노드 그룹**: 워커 노드들
- **IAM 역할**: EKS 클러스터 및 노드 권한

### **Kubernetes 리소스**
- **aws-auth ConfigMap**: GitHub Actions 접근 권한

## ⚙️ 설정 변수

주요 변수들을 `variables.tf`에서 수정할 수 있습니다:

```hcl
variable "aws_region" {
  default = "us-east-1"
}

variable "project_name" {
  default = "github-actions-terraform"
}

variable "eks_cluster_name" {
  default = "sm-eks"
}

variable "github_actions_user" {
  default = "sm-user"
}
```

## 🔧 문제 해결

### EKS 접근 권한 오류
```bash
# 수동으로 권한 추가
eksctl create iamidentitymapping \
  --cluster sm-eks \
  --region us-east-1 \
  --arn arn:aws:iam::ACCOUNT-ID:user/sm-user \
  --group system:masters \
  --username github-actions
```

### Terraform State 문제
```bash
# State 새로고침
terraform refresh

# State 파일 재생성 (주의!)
terraform import aws_eks_cluster.main sm-eks
```

## 🗑️ 정리

### 전체 삭제
```bash
terraform destroy
```

### EKS만 삭제
```bash
eksctl delete cluster --name sm-eks --region us-east-1
```

## 📋 출력값

배포 완료 후 다음 정보들을 확인할 수 있습니다:

- **ECR Repository URL**: Docker 이미지 푸시 주소
- **EKS Cluster Endpoint**: Kubernetes API 서버 주소
- **EKS Cluster Name**: 클러스터 이름
- **IAM Role ARNs**: 생성된 역할들의 ARN

## 🔗 연관 파일

- **GitHub Actions**: `.github/workflows/terraform-infrastructure.yml`
- **애플리케이션**: `app/` 디렉토리
- **Kubernetes**: `k8s/` 디렉토리