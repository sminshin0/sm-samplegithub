package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// 포트 설정 (환경변수 또는 기본값 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 라우트 설정
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Hello World Go App</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: #00ADD8; }
        .info { background: #f0f0f0; padding: 20px; border-radius: 10px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎉 Hello World from Go!</h1>
        <div class="info">
            <p><strong>Version:</strong> 1.0.0</p>
            <p><strong>Built with:</strong> Go + Docker + GitHub Actions</p>
            <p><strong>Deployed to:</strong> AWS ECR</p>
        </div>
        <p>✅ Application is running successfully!</p>
    </div>
</body>
</html>`)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "healthy", "message": "Go web server is running"}`)
}