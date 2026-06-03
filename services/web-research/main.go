package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := getEnv("PORT", "8083")
	log.Printf("Iniciando Web Research Service na porta %s", port)

	handler := &Handler{
		researcher: NewResearcher(&http.Client{Timeout: 8 * time.Second}),
	}

	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)
	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/research", handler.ResearchHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	log.Printf("Web Research Service escutando em 0.0.0.0:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
