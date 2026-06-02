package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler struct para guardar dependências
type Handler struct {
	db *Database
}

// SearchCatalogHandler GET /catalog/search?pn=XXXXX
func (h *Handler) SearchCatalogHandler(w http.ResponseWriter, r *http.Request) {
	partNumber := r.URL.Query().Get("pn")

	if partNumber == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Part Number é obrigatório",
		})
		return
	}

	log.Printf("Buscando PN: %s", partNumber)

	item, err := h.db.SearchCatalog(partNumber)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SearchResponse{
			Found: false,
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(SearchResponse{Found: false})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchResponse{
		Found: true,
		Item:  item,
	})
}

// AddInventoryHandler POST /inventory/in
func (h *Handler) AddInventoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req InventoryInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao decodificar requisição",
		})
		return
	}

	log.Printf("Adicionando item ao estoque: PN=%s, SN=%s, QTY=%d",
		req.PartNumber, req.SerialNumber, req.Quantity)

	if req.PartNumber == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(InventoryResponse{
			Success: false,
			Error:   "Part Number é obrigatório",
		})
		return
	}

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	inventory, err := h.db.AddInventory(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InventoryResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(InventoryResponse{
		Success: true,
		Data:    inventory,
	})
}

// GetInventoryHandler GET /inventory/items/{id}
func (h *Handler) GetInventoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "ID inválido",
		})
		return
	}

	item, err := h.db.GetInventory(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if item == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Item não encontrado",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// ListInventoryHandler GET /inventory/items?limit=50&offset=0
func (h *Handler) ListInventoryHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > 500 {
				parsed = 500
			}
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	items, total, err := h.db.ListInventory(limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListResponse{
		Total: total,
		Items: items,
	})
}

// HealthHandler GET /health
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "inventory-service",
		"version": "1.0.0",
	})
}

// RegisterRoutes registra as rotas
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.HealthHandler).Methods("GET")
	router.HandleFunc("/catalog/search", h.SearchCatalogHandler).Methods("GET")
	router.HandleFunc("/inventory/in", h.AddInventoryHandler).Methods("POST")
	router.HandleFunc("/inventory/items", h.ListInventoryHandler).Methods("GET")
	router.HandleFunc("/inventory/items/{id}", h.GetInventoryHandler).Methods("GET")
}
