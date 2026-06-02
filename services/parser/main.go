package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := getEnv("PORT", "8082")
	log.Printf("Iniciando Parser Service na porta %s", port)

	handler := &Handler{}
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/parse", handler.ParseHandler).Methods(http.MethodPost)
	router.HandleFunc("/parse", handler.ParseHandler).Methods(http.MethodOptions)

	log.Printf("Parser Service escutando em 0.0.0.0:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
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

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
