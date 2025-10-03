# Kubernetes Manifests

Go 웹 애플리케이션을 위한 Kubernetes 배포 매니페스트입니다.

## 📁 구조

```
k8s/
├── deployment.yaml    # 애플리케이션 배포
├── service.yaml       # 서비스 (LoadBalancer + NodePort)
├── configmap.yaml     # 설정 관리
└── README.md          # 이 파일
```

## 🚀 배포 방법

### 수동 배포
```bash
# 모든 리소스 배포
kubectl apply -f k8s/

# 개별 리소스 배포
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 자동 배포
- `k8s/` 폴더 내 파일 변경 시 GitHub Actions가 자동으로 EKS에 배포

## 🔧 리소스 설명

### Deployment
- **Replicas**: 2개 (고가용성)
- **Resources**: CPU 50m-100m, Memory 64Mi-128Mi
- **Health Checks**: Liveness + Readiness Probe
- **Image**: ECR에서 자동으로 최신 이미지 사용

### Service
- **LoadBalancer**: 외부 접근용 (포트 80 → 8080)
- **NodePort**: 개발/테스트용 (포트 30080)

### ConfigMap
- 환경 변수 및 설정 관리
- 애플리케이션 재빌드 없이 설정 변경 가능

## 🌐 접근 방법

### LoadBalancer 사용
```bash
# 외부 IP 확인
kubectl get service go-hello-world-service

# 브라우저에서 접근
# http://<EXTERNAL-IP>
```

### NodePort 사용 (개발용)
```bash
# 노드 IP 확인
kubectl get nodes -o wide

# 브라우저에서 접근
# http://<NODE-IP>:30080
```

## 📊 모니터링

```bash
# Pod 상태 확인
kubectl get pods -l app=go-hello-world

# 로그 확인
kubectl logs -l app=go-hello-world -f

# 서비스 상태 확인
kubectl get services

# 배포 상태 확인
kubectl rollout status deployment/go-hello-world
```

## 🔄 업데이트

```bash
# 이미지 업데이트 (수동)
kubectl set image deployment/go-hello-world go-hello-world=<NEW_IMAGE>

# 배포 롤백
kubectl rollout undo deployment/go-hello-world

# 배포 히스토리
kubectl rollout history deployment/go-hello-world
```