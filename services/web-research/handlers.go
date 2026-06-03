package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Handler struct {
	researcher *Researcher
}

func (h *Handler) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "web-research-service",
		"version": "1.0.0",
	})
}

func (h *Handler) ResearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req ResearchRequest
	if r.Method == http.MethodGet {
		req.PartNumber = strings.TrimSpace(r.URL.Query().Get("pn"))
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ResearchResponse{Success: false, Error: "Erro ao decodificar requisição"})
			return
		}
	}

	req.PartNumber = strings.TrimSpace(strings.ToUpper(req.PartNumber))
	if req.PartNumber == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ResearchResponse{Success: false, Error: "Part Number é obrigatório"})
		return
	}

	response := h.researcher.Research(r.Context(), req)
	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusBadGateway)
	}
	_ = json.NewEncoder(w).Encode(response)
}
