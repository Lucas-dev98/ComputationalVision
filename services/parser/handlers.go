package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct{}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "parser-service",
		"version": "1.0.0",
	})
}

func (h *Handler) ParseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ParseResponse{Success: false, Error: "Erro ao decodificar requisição"})
		return
	}

	response := ParseLines(req.Text)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
