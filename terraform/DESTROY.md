# 🗑️ Terraform 인프라 삭제 가이드

## ⚠️ 주의사항
- **모든 데이터가 영구적으로 삭제됩니다**
- **ECR의 Docker 이미지들도 삭제됩니다**
- **S3 버킷의 Terraform State도 삭제됩니다**
- **되돌릴 수 없으니 신중하게 결정하세요**

## 🚀 방법 1: GitHub Actions 사용 (권장)

1. **GitHub 리포지토리** → **Actions** 탭
2. **"Terraform Destroy"** 워크플로우 선택
3. **"Run workflow"** 클릭
4. **"confirm_destroy"**를 **"yes"**로 선택
5. **"Run workflow"** 실행
6. 진행 상황을 모니터링

## 💻 방법 2: 로컬에서 직접 삭제

### 사전 준비
```bash
# AWS CLI 설정 확인
aws sts get-caller-identity

# terraform 폴더로 이동
cd terraform

# 현재 상태 확인
terraform plan
```

### 1단계: 사전 정리
```bash
# ECR 이미지들 삭제
aws ecr list-images --repository-name github-actions-terraform-app --region us-east-1 --query 'imageIds[*]' --output json > images.json
aws ecr batch-delete-image --repository-name github-actions-terraform-app --region us-east-1 --image-ids file://images.json

# EKS 클러스터 연결 및 서비스 정리
aws eks update-kubeconfig --region us-east-1 --name sm-eks
kubectl delete service --all -n default --timeout=300s
```

### 2단계: Terraform Destroy
```bash
# 삭제 계획 확인
terraform plan -destroy

# 실제 삭제 실행 (주의!)
terraform destroy

# 확인 메시지에서 'yes' 입력
```

### 3단계: 후속 정리
```bash
# S3 버킷 강제 삭제 (버전이 있는 경우)
aws s3 rm s3://terraform-state-us-east-1-0ss4kx0a --recursive
aws s3 rb s3://terraform-state-us-east-1-0ss4kx0a --force

# DynamoDB 테이블 확인
aws dynamodb describe-table --table-name my-terraform-project-terraform-state-lock --region us-east-1
```

## 🔍 삭제 확인 방법

### AWS 콘솔에서 확인
1. **EKS**: 클러스터가 삭제되었는지 확인
2. **ECR**: 리포지토리가 삭제되었는지 확인
3. **IAM**: 생성된 역할들이 삭제되었는지 확인
4. **EC2**: 보안 그룹이 삭제되었는지 확인
5. **S3**: terraform state 버킷이 삭제되었는지 확인
6. **DynamoDB**: lock 테이블이 삭제되었는지 확인

### AWS CLI로 확인
```bash
# EKS 클러스터 확인
aws eks list-clusters --region us-east-1

# ECR 리포지토리 확인
aws ecr describe-repositories --region us-east-1

# S3 버킷 확인
aws s3 ls | grep terraform-state

# DynamoDB 테이블 확인
aws dynamodb list-tables --region us-east-1 | grep terraform
```

## 💰 비용 절약 효과

삭제 후 절약되는 월 예상 비용:
- **EKS 클러스터**: ~$73/월
- **EKS 노드 그룹** (t3.medium x2): ~$60/월
- **ECR 스토리지**: ~$1-5/월
- **기타 리소스**: ~$5/월
- **총 절약**: ~$140/월

## 🔄 재생성 방법

필요시 다시 생성하려면:
```bash
# terraform 폴더에서
terraform init
terraform plan
terraform apply
```

## 🆘 문제 해결

### Terraform Destroy 실패 시
```bash
# 특정 리소스만 삭제
terraform destroy -target=aws_eks_cluster.main

# 강제 삭제 (주의!)
terraform state rm aws_eks_cluster.main
```

### 수동 정리가 필요한 경우
```bash
# EKS 클러스터 수동 삭제
aws eks delete-cluster --name sm-eks --region us-east-1

# ECR 리포지토리 수동 삭제
aws ecr delete-repository --repository-name github-actions-terraform-app --region us-east-1 --force
```

## 📞 지원

문제가 발생하면:
1. GitHub Issues에 문제 보고
2. AWS 콘솔에서 수동 정리
3. AWS Support 문의 (필요시)