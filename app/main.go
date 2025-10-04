package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type InfoResponse struct {
	AppName     string            `json:"app_name"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
	Port        string            `json:"port"`
	Uptime      time.Duration     `json:"uptime"`
	Headers     map[string]string `json:"headers"`
}

var startTime = time.Now()

func main() {
	// 로거 설정
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// 포트 설정 (환경변수 또는 기본값 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Gorilla Mux 라우터 생성
	r := mux.NewRouter()

	// 라우트 설정
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/info", infoHandler).Methods("GET")
	r.HandleFunc("/api/time", timeHandler).Methods("GET")
	r.HandleFunc("/api/deployment", deploymentHandler).Methods("GET")

	// 미들웨어 추가
	r.Use(loggingMiddleware)

	logrus.WithFields(logrus.Fields{
		"port":    port,
		"version": "1.0.0",
	}).Info("🚀 Server starting")

	// 서버 시작
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
			"ip":       r.RemoteAddr,
		}).Info("Request processed")
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Go Web Application - EKS Deployment</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 30px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; }
        .container { max-width: 900px; margin: 0 auto; background: rgba(255,255,255,0.95); padding: 40px; border-radius: 15px; box-shadow: 0 8px 32px rgba(0,0,0,0.3); color: #333; }
        h1 { color: #00ADD8; margin-bottom: 30px; font-size: 2.5em; }
        .deployment-info { background: linear-gradient(45deg, #28a745, #20c997); color: white; padding: 20px; border-radius: 10px; margin: 20px 0; }
        .info { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; text-align: left; }
        .api-links { margin: 20px 0; }
        .api-links a { display: inline-block; margin: 5px 10px; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 25px; transition: all 0.3s; }
        .api-links a:hover { background: #0056b3; transform: translateY(-2px); }
        .stats { display: flex; justify-content: space-around; margin: 30px 0; }
        .stat { text-align: center; }
        .stat-value { font-size: 2.2em; font-weight: bold; color: #28a745; }
        .new-feature { background: #ff6b6b; color: white; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .timestamp { font-size: 0.9em; color: #666; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Go Web Application</h1>
        
        <div class="deployment-info">
            <h3>🎯 Successfully Deployed to AWS EKS!</h3>
            <p>This application is running on Kubernetes cluster with automated CI/CD</p>
            <p><strong>Deployment Time:</strong> %s</p>
        </div>
        
        <div class="stats">
            <div class="stat">
                <div class="stat-value">%s</div>
                <div>Uptime</div>
            </div>
            <div class="stat">
                <div class="stat-value">v2.0.0</div>
                <div>Version</div>
            </div>
            <div class="stat">
                <div class="stat-value">🎉</div>
                <div>EKS Ready</div>
            </div>
        </div>

        <div class="new-feature">
            <h3>🆕 New Features Added!</h3>
            <ul style="text-align: left; margin: 10px 0;">
                <li>✨ Enhanced UI with gradient design</li>
                <li>🚀 EKS deployment information</li>
                <li>📊 Real-time deployment status</li>
                <li>🔄 Automated GitHub Actions CI/CD</li>
            </ul>
        </div>

        <div class="info">
            <h3>🛠️ Infrastructure Stack</h3>
            <ul>
                <li><strong>Language:</strong> Go 1.21 with Gorilla Mux</li>
                <li><strong>Container:</strong> Docker (Multi-stage build)</li>
                <li><strong>Registry:</strong> AWS ECR</li>
                <li><strong>Orchestration:</strong> AWS EKS (Kubernetes)</li>
                <li><strong>Infrastructure:</strong> Terraform (IaC)</li>
                <li><strong>CI/CD:</strong> GitHub Actions</li>
                <li><strong>Monitoring:</strong> Structured logging with Logrus</li>
            </ul>
        </div>

        <div class="api-links">
            <h3>🔗 API Endpoints</h3>
            <a href="/health">Health Check</a>
            <a href="/info">App Info</a>
            <a href="/api/time">Current Time</a>
            <a href="/api/deployment">Deployment Info</a>
        </div>

        <div class="timestamp">
            <p>🕒 Last updated: %s | 🌍 Running on EKS cluster</p>
        </div>
        
        <p style="font-size: 1.2em; color: #28a745;">✅ Application successfully deployed and running on AWS EKS!</p>
    </div>
</body>
</html>`, time.Now().Format("2006-01-02 15:04:05 UTC"), uptime.Round(time.Second), time.Now().Format("2006-01-02 15:04:05"))

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Message:   "Go web server is running",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string]string)
	for name, values := range r.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	response := InfoResponse{
		AppName:     "github-actions-terraform-app",
		Version:     "1.0.0",
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		Uptime:      time.Since(startTime),
		Headers:     headers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"current_time": time.Now(),
		"unix_time":    time.Now().Unix(),
		"timezone":     time.Now().Format("MST"),
		"iso8601":      time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deploymentHandler(w http.ResponseWriter, r *http.Request) {
	deployment := map[string]interface{}{
		"application": map[string]string{
			"name":    "github-actions-terraform-app",
			"version": "2.0.0",
			"status":  "running",
		},
		"infrastructure": map[string]string{
			"platform":    "AWS EKS",
			"region":      getEnv("AWS_REGION", "us-east-1"),
			"environment": getEnv("ENVIRONMENT", "production"),
		},
		"deployment": map[string]interface{}{
			"method":     "GitHub Actions + Terraform",
			"container":  "Docker (ECR)",
			"started_at": startTime,
			"uptime":     time.Since(startTime),
		},
		"kubernetes": map[string]string{
			"namespace": getEnv("KUBERNETES_NAMESPACE", "default"),
			"pod_name":  getEnv("HOSTNAME", "unknown"),
		},
		"build_info": map[string]string{
			"go_version": "1.21",
			"git_commit": getEnv("GIT_COMMIT", "unknown"),
			"build_time": getEnv("BUILD_TIME", time.Now().Format(time.RFC3339)),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}