#!/bin/bash

# EKS aws-auth 수동 설정 스크립트
set -e

CLUSTER_NAME="sm-eks"
REGION="us-east-1"
GITHUB_USER="sm-user"
ACCOUNT_ID="727158410149"

echo "🔧 EKS aws-auth ConfigMap 설정 중..."

# EKS 클러스터 연결
echo "📡 EKS 클러스터 연결 중..."
aws eks update-kubeconfig --region $REGION --name $CLUSTER_NAME

# 현재 aws-auth 확인
echo "🔍 현재 aws-auth ConfigMap 확인 중..."
if kubectl get configmap aws-auth -n kube-system >/dev/null 2>&1; then
    echo "✅ aws-auth ConfigMap 존재"
    kubectl get configmap aws-auth -n kube-system -o yaml > aws-auth-backup.yaml
    echo "📋 백업 파일 생성: aws-auth-backup.yaml"
else
    echo "❌ aws-auth ConfigMap이 없습니다. 새로 생성합니다."
fi

# eksctl로 IAM 사용자 추가
echo "👤 GitHub Actions 사용자 권한 추가 중..."
eksctl create iamidentitymapping \
    --cluster $CLUSTER_NAME \
    --region $REGION \
    --arn "arn:aws:iam::$ACCOUNT_ID:user/$GITHUB_USER" \
    --group system:masters \
    --username github-actions

# 권한 확인
echo "🧪 권한 테스트 중..."
kubectl auth can-i get pods --as=github-actions

echo "✅ EKS aws-auth 설정 완료!"
echo ""
echo "🎉 이제 GitHub Actions에서 EKS 클러스터에 접근할 수 있습니다."
echo ""
echo "📋 설정된 내용:"
echo "   - 클러스터: $CLUSTER_NAME"
echo "   - 사용자: arn:aws:iam::$ACCOUNT_ID:user/$GITHUB_USER"
echo "   - 권한: system:masters"
echo ""
echo "🚀 다음 단계: GitHub Actions 워크플로우를 실행해보세요!"