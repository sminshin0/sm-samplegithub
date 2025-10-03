package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type HelloResponse struct {
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
    Version   string `json:"version"`
    ClientIP  string `json:"client_ip"`
}

type HealthResponse struct {
    Status string `json:"status"`
    Uptime string `json:"uptime"`
}

var startTime = time.Now()

func helloAPIHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    response := HelloResponse{
        Message:   "Hello, World from Go API! 🚀",
        Timestamp: time.Now().Format(time.RFC3339),
        Version:   "1.0.0",
        ClientIP:  r.RemoteAddr,
    }
    
    json.NewEncoder(w).Encode(response)
}

func healthAPIHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    uptime := time.Since(startTime).Round(time.Second)
    
    response := HealthResponse{
        Status: "healthy",
        Uptime: uptime.String(),
    }
    
    json.NewEncoder(w).Encode(response)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func main() {
    http.HandleFunc("/api/hello", corsMiddleware(helloAPIHandler))
    http.HandleFunc("/api/health", corsMiddleware(healthAPIHandler))
    
    port := ":8080"
    fmt.Printf("🚀 Go API Server starting on port %s\n", port)
    fmt.Println("📍 API Endpoints:")
    fmt.Println("   GET http://localhost:8080/api/hello")
    fmt.Println("   GET http://localhost:8080/api/health")
    fmt.Println("")
    fmt.Println("💡 Try: curl http://localhost:8080/api/hello")
    
    log.Fatal(http.ListenAndServe(port, nil))
}