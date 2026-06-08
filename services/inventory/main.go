package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Carregar configurações
	config := LoadConfig()
	log.Printf("Iniciando Inventory Service na porta %s", config.Port)

	// Conectar ao banco
	db, err := NewDatabase(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}
	defer db.Close()

	log.Println("Conectado ao banco com sucesso")

	// Criar handler
	handler := &Handler{db: db}

	// Criar router
	router := mux.NewRouter()

	// Registrar rotas
	handler.RegisterRoutes(router)

	// Middleware para logs
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	router.Use(newRateLimitMiddleware(loadRateLimitConfig()))

	// Iniciar servidor
	log.Printf("Servidor escutando em 0.0.0.0:%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, router); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

// loggingMiddleware faz log de requisições
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adiciona headers CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
