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
    <title>Go Web Application</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; background: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #00ADD8; margin-bottom: 30px; }
        .info { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; text-align: left; }
        .api-links { margin: 20px 0; }
        .api-links a { display: inline-block; margin: 5px 10px; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; }
        .api-links a:hover { background: #0056b3; }
        .stats { display: flex; justify-content: space-around; margin: 20px 0; }
        .stat { text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; color: #28a745; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎉 Go Web Application</h1>
        
        <div class="stats">
            <div class="stat">
                <div class="stat-value">%s</div>
                <div>Uptime</div>
            </div>
            <div class="stat">
                <div class="stat-value">v1.0.0</div>
                <div>Version</div>
            </div>
            <div class="stat">
                <div class="stat-value">✅</div>
                <div>Status</div>
            </div>
        </div>

        <div class="info">
            <h3>🛠️ Tech Stack</h3>
            <ul>
                <li><strong>Language:</strong> Go 1.21</li>
                <li><strong>Router:</strong> Gorilla Mux</li>
                <li><strong>Logging:</strong> Logrus</li>
                <li><strong>Container:</strong> Docker</li>
                <li><strong>Deployment:</strong> AWS EKS + ECR</li>
                <li><strong>CI/CD:</strong> GitHub Actions</li>
            </ul>
        </div>

        <div class="api-links">
            <h3>🔗 API Endpoints</h3>
            <a href="/health">Health Check</a>
            <a href="/info">App Info</a>
            <a href="/api/time">Current Time</a>
        </div>

        <p>✅ Application is running successfully!</p>
    </div>
</body>
</html>`, uptime.Round(time.Second))

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}