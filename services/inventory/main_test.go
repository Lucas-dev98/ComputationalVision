package main
package main

import (
	"testing"
	"os"
)

func TestDatabaseConnection(t *testing.T) {
	// Este é um teste de exemplo
	// Em produção, usar um banco de testes
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL não definida, skipping database tests")
	}
	
	db, err := NewDatabase(dsn)
	if err != nil {
		t.Fatalf("Falha ao conectar ao banco: %v", err)
	}
	defer db.Close()
	
	t.Log("✓ Conexão com banco de dados OK")
}

func TestSearchCatalog(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL não definida, skipping database tests")
	}
	
	db, err := NewDatabase(dsn)
	if err != nil {
		t.Fatalf("Falha ao conectar ao banco: %v", err)
	}
	defer db.Close()
	
	// Buscar um item que deve existir (dos dados de teste)
	item, err := db.SearchCatalog("M393A4K40DB3-CWE")
	if err != nil {
		t.Fatalf("Erro ao buscar catálogo: %v", err)
	}
	
	if item == nil {
		t.Error("Item não encontrado no catálogo")
		return
	}
	
	if item.PartNumber != "M393A4K40DB3-CWE" {
		t.Errorf("PN incorreto: esperado M393A4K40DB3-CWE, got %s", item.PartNumber)
	}
	
	t.Log("✓ Busca em catálogo OK")
}

func TestSearchCatalogNotFound(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL não definida, skipping database tests")
	}
	
	db, err := NewDatabase(dsn)
	if err != nil {
		t.Fatalf("Falha ao conectar ao banco: %v", err)
	}
	defer db.Close()
	
	// Buscar um item que não existe
	item, err := db.SearchCatalog("INVALID_PN_XXXX")
	if err != nil {
		t.Fatalf("Erro ao buscar catálogo: %v", err)
	}
	
	if item != nil {
		t.Error("Item não deveria ter sido encontrado")
	}
	
	t.Log("✓ Busca negativa em catálogo OK")
}
