package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World! 🌍\n")
    fmt.Fprintf(w, "Current time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Fprintf(w, "Request from: %s\n", r.RemoteAddr)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/health", healthHandler)
    
    port := ":8080"
    fmt.Printf("🚀 Server starting on port %s\n", port)
    fmt.Println("📍 Endpoints:")
    fmt.Println("   http://localhost:8080/")
    fmt.Println("   http://localhost:8080/health")
    
    log.Fatal(http.ListenAndServe(port, nil))
}