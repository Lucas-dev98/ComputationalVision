package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Este é um teste de exemplo
	// Em produção, usar um banco de testes

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL não definida, skipping database tests")
	}

	db, err := NewDatabase("postgres", dsn)
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

	db, err := NewDatabase("postgres", dsn)
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

	db, err := NewDatabase("postgres", dsn)
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

func TestSQLiteEmbeddedMode(t *testing.T) {
	tempDB := filepath.Join(t.TempDir(), "inventory-test.db")
	db, err := NewDatabase("sqlite", "file:"+tempDB+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("Falha ao iniciar SQLite embutido: %v", err)
	}
	defer db.Close()

	item, err := db.SearchCatalog("M393A4K40DB3-CWE")
	if err != nil {
		t.Fatalf("Erro na busca em catálogo SQLite: %v", err)
	}
	if item == nil {
		t.Fatal("Seed do catálogo não foi criado no SQLite")
	}

	inv, err := db.AddInventory(&InventoryInRequest{
		PartNumber:   "M393A4K40DB3-CWE",
		SerialNumber: "SN-SQLITE-001",
		Quantity:     1,
		Location:     "LAB-LOCAL",
		Reason:       "Teste local",
		UserID:       1,
	})
	if err != nil {
		t.Fatalf("Erro ao adicionar inventário no SQLite: %v", err)
	}
	if inv.ID == 0 {
		t.Fatal("ID do inventário inválido")
	}
}
