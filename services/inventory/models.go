package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// Database gerencia a conexão com PostgreSQL
type Database struct {
	conn *sql.DB
}

// CatalogItem representa um item no catálogo
type CatalogItem struct {
	ID                    int64     `json:"id"`
	PartNumber            string    `json:"part_number"`
	SerialPattern         string    `json:"serial_pattern"`
	Manufacturer          string    `json:"manufacturer"`
	Category              string    `json:"category"`
	NormalizedDescription string    `json:"normalized_description"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// InventoryItem representa um item em estoque
type InventoryItem struct {
	ID           int64        `json:"id"`
	CatalogID    int64        `json:"catalog_id"`
	SerialNumber string       `json:"serial_number"`
	Quantity     int          `json:"quantity"`
	Location     string       `json:"location"`
	Status       string       `json:"status"`
	ReceivedAt   time.Time    `json:"received_at"`
	LastUpdated  time.Time    `json:"last_updated"`
	Catalog      *CatalogItem `json:"catalog,omitempty"`
}

// Movement representa uma movimentação de estoque
type Movement struct {
	ID          int64     `json:"id"`
	InventoryID int64     `json:"inventory_id"`
	Quantity    int       `json:"quantity"`
	Operation   string    `json:"operation"`
	Reason      string    `json:"reason"`
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// FeedbackSample representa uma amostra de correção humana para aprendizado contínuo.
type FeedbackSample struct {
	ID                    int64     `json:"id"`
	PartNumberPredicted   string    `json:"part_number_predicted,omitempty"`
	PartNumberFinal       string    `json:"part_number_final,omitempty"`
	SerialNumberPredicted string    `json:"serial_number_predicted,omitempty"`
	SerialNumberFinal     string    `json:"serial_number_final,omitempty"`
	ManufacturerPredicted string    `json:"manufacturer_predicted,omitempty"`
	ManufacturerFinal     string    `json:"manufacturer_final,omitempty"`
	CategoryPredicted     string    `json:"category_predicted,omitempty"`
	CategoryFinal         string    `json:"category_final,omitempty"`
	CorrectionApplied     bool      `json:"correction_applied"`
	Confidence            float64   `json:"confidence,omitempty"`
	ImageData             string    `json:"image_data,omitempty"`
	OCRText               []string  `json:"ocr_text,omitempty"`
	MetaJSON              string    `json:"meta_json,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
}

type FeedbackRequest struct {
	PartNumberPredicted   string   `json:"part_number_predicted"`
	PartNumberFinal       string   `json:"part_number_final"`
	SerialNumberPredicted string   `json:"serial_number_predicted"`
	SerialNumberFinal     string   `json:"serial_number_final"`
	ManufacturerPredicted string   `json:"manufacturer_predicted"`
	ManufacturerFinal     string   `json:"manufacturer_final"`
	CategoryPredicted     string   `json:"category_predicted"`
	CategoryFinal         string   `json:"category_final"`
	CorrectionApplied     bool     `json:"correction_applied"`
	Confidence            float64  `json:"confidence"`
	ImageData             string   `json:"image_data"`
	OCRText               []string `json:"ocr_text"`
	MetaJSON              string   `json:"meta_json"`
}

type FeedbackResponse struct {
	Success bool            `json:"success"`
	Data    *FeedbackSample `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type FeedbackListResponse struct {
	Total int              `json:"total"`
	Items []FeedbackSample `json:"items"`
}

// InventoryInRequest requisição para entrada de estoque
type InventoryInRequest struct {
	PartNumber   string `json:"part_number"`
	SerialNumber string `json:"serial_number"`
	Quantity     int    `json:"quantity"`
	Location     string `json:"location"`
	Reason       string `json:"reason"`
	UserID       int64  `json:"user_id"`
}

// InventoryOutRequest requisição para saída de estoque
type InventoryOutRequest struct {
	InventoryID int64  `json:"inventory_id"`
	Quantity    int    `json:"quantity"`
	Reason      string `json:"reason"`
	UserID      int64  `json:"user_id"`
}

// SearchResponse resposta de busca
type SearchResponse struct {
	Found bool         `json:"found"`
	Item  *CatalogItem `json:"item,omitempty"`
	Error string       `json:"error,omitempty"`
}

// InventoryResponse resposta de operação de estoque
type InventoryResponse struct {
	Success bool           `json:"success"`
	Data    *InventoryItem `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// ListResponse resposta de listagem
type ListResponse struct {
	Total int             `json:"total"`
	Items []InventoryItem `json:"items"`
}

// Config armazena configurações da aplicação
type Config struct {
	DatabaseDriver string
	DatabaseURL    string
	RedisURL       string
	Port           string
	LogLevel       string
}

// LoadConfig carrega configurações do ambiente
func LoadConfig() *Config {
	godotenv.Load()

	driver := strings.ToLower(getEnv("DATABASE_DRIVER", "postgres"))
	defaultDSN := "postgres://inventory:inventory_dev@localhost:5432/inventory_db?sslmode=disable"
	if driver == "sqlite" {
		defaultDSN = "file:./inventory-dev.db?_pragma=foreign_keys(1)"
	}

	return &Config{
		DatabaseDriver: driver,
		DatabaseURL:    getEnv("DATABASE_URL", defaultDSN),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		Port:           getEnv("PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv obtém variável de ambiente com padrão
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// NewDatabase cria nova conexão com banco
func NewDatabase(driver, dsn string) (*Database, error) {
	if driver == "" {
		driver = "postgres"
	}

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		log.Printf("Erro ao conectar ao banco: %v", err)
		return nil, err
	}

	// Testar conexão
	if err := conn.Ping(); err != nil {
		log.Printf("Erro ao verificar conexão: %v", err)
		return nil, err
	}

	// Configurar pool
	if driver == "sqlite" {
		conn.SetMaxOpenConns(1)
		conn.SetMaxIdleConns(1)
	} else {
		conn.SetMaxOpenConns(25)
		conn.SetMaxIdleConns(5)
	}
	conn.SetConnMaxLifetime(5 * time.Minute)

	if driver == "sqlite" {
		if err := bootstrapSQLite(conn); err != nil {
			return nil, err
		}
		log.Println("Conectado ao SQLite (modo local dev) com sucesso")
	} else {
		log.Println("Conectado ao PostgreSQL com sucesso")
	}

	return &Database{conn: conn}, nil
}

// Close fecha conexão
func (db *Database) Close() error {
	return db.conn.Close()
}

// SearchCatalog busca um item no catálogo
func (db *Database) SearchCatalog(partNumber string) (*CatalogItem, error) {
	item := &CatalogItem{}

	err := db.conn.QueryRow(`
		SELECT id, part_number,
		       COALESCE(serial_pattern, ''),
		       COALESCE(manufacturer, ''),
		       COALESCE(category, ''),
		       COALESCE(normalized_description, ''),
		       created_at, updated_at
		FROM catalog
		WHERE part_number = $1
		LIMIT 1
	`, partNumber).Scan(
		&item.ID,
		&item.PartNumber,
		&item.SerialPattern,
		&item.Manufacturer,
		&item.Category,
		&item.NormalizedDescription,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Printf("Erro ao buscar no catálogo: %v", err)
		return nil, err
	}

	return item, nil
}

// AddInventory adiciona item ao estoque
func (db *Database) AddInventory(req *InventoryInRequest) (*InventoryItem, error) {
	// Buscar item no catálogo
	catalogItem, err := db.SearchCatalog(req.PartNumber)
	if err != nil {
		return nil, err
	}

	if catalogItem == nil {
		log.Printf("Part Number não encontrado: %s", req.PartNumber)
		return nil, sql.ErrNoRows
	}

	// Inserir em estoque
	inventory := &InventoryItem{}
	err = db.conn.QueryRow(`
		INSERT INTO inventory (catalog_id, serial_number, quantity, location, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, catalog_id, serial_number, quantity, location, status, received_at, last_updated
	`,
		catalogItem.ID,
		req.SerialNumber,
		req.Quantity,
		req.Location,
		"active",
	).Scan(
		&inventory.ID,
		&inventory.CatalogID,
		&inventory.SerialNumber,
		&inventory.Quantity,
		&inventory.Location,
		&inventory.Status,
		&inventory.ReceivedAt,
		&inventory.LastUpdated,
	)

	if err != nil {
		log.Printf("Erro ao inserir em estoque: %v", err)
		return nil, err
	}

	// Registrar movimento
	_, err = db.conn.Exec(`
		INSERT INTO movements (inventory_id, quantity, operation, reason, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`,
		inventory.ID,
		req.Quantity,
		"IN",
		req.Reason,
		req.UserID,
	)

	if err != nil {
		log.Printf("Erro ao registrar movimento: %v", err)
	}

	inventory.Catalog = catalogItem
	return inventory, nil
}

// GetInventory obtém item do estoque com seu catálogo
func (db *Database) GetInventory(id int64) (*InventoryItem, error) {
	inventory := &InventoryItem{}
	catalogItem := &CatalogItem{}

	err := db.conn.QueryRow(`
		SELECT i.id, i.catalog_id,
		       COALESCE(i.serial_number, ''),
		       i.quantity,
		       COALESCE(i.location, ''),
		       i.status, i.received_at, i.last_updated,
		       c.id, c.part_number,
		       COALESCE(c.serial_pattern, ''),
		       COALESCE(c.manufacturer, ''),
		       COALESCE(c.category, ''),
		       COALESCE(c.normalized_description, ''),
		       c.created_at, c.updated_at
		FROM inventory i
		JOIN catalog c ON i.catalog_id = c.id
		WHERE i.id = $1
	`, id).Scan(
		&inventory.ID,
		&inventory.CatalogID,
		&inventory.SerialNumber,
		&inventory.Quantity,
		&inventory.Location,
		&inventory.Status,
		&inventory.ReceivedAt,
		&inventory.LastUpdated,
		&catalogItem.ID,
		&catalogItem.PartNumber,
		&catalogItem.SerialPattern,
		&catalogItem.Manufacturer,
		&catalogItem.Category,
		&catalogItem.NormalizedDescription,
		&catalogItem.CreatedAt,
		&catalogItem.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Printf("Erro ao buscar item: %v", err)
		return nil, err
	}

	inventory.Catalog = catalogItem
	return inventory, nil
}

// ListInventory lista itens em estoque
func (db *Database) ListInventory(limit, offset int) ([]InventoryItem, int, error) {
	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM inventory WHERE status = 'active'").Scan(&total)
	if err != nil {
		log.Printf("Erro ao contar itens: %v", err)
		return nil, 0, err
	}

	rows, err := db.conn.Query(`
		SELECT i.id, i.catalog_id,
		       COALESCE(i.serial_number, ''),
		       i.quantity,
		       COALESCE(i.location, ''),
		       i.status, i.received_at, i.last_updated,
		       c.id, c.part_number,
		       COALESCE(c.serial_pattern, ''),
		       COALESCE(c.manufacturer, ''),
		       COALESCE(c.category, ''),
		       COALESCE(c.normalized_description, ''),
		       c.created_at, c.updated_at
		FROM inventory i
		JOIN catalog c ON i.catalog_id = c.id
		WHERE i.status = 'active'
		ORDER BY i.received_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Erro ao listar itens: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	items := []InventoryItem{}
	for rows.Next() {
		inventory := InventoryItem{}
		catalogItem := CatalogItem{}

		err := rows.Scan(
			&inventory.ID,
			&inventory.CatalogID,
			&inventory.SerialNumber,
			&inventory.Quantity,
			&inventory.Location,
			&inventory.Status,
			&inventory.ReceivedAt,
			&inventory.LastUpdated,
			&catalogItem.ID,
			&catalogItem.PartNumber,
			&catalogItem.SerialPattern,
			&catalogItem.Manufacturer,
			&catalogItem.Category,
			&catalogItem.NormalizedDescription,
			&catalogItem.CreatedAt,
			&catalogItem.UpdatedAt,
		)

		if err != nil {
			log.Printf("Erro ao escanear item: %v", err)
			continue
		}

		inventory.Catalog = &catalogItem
		items = append(items, inventory)
	}

	return items, total, rows.Err()
}

func (db *Database) AddFeedback(req *FeedbackRequest) (*FeedbackSample, error) {
	ocrTextJSON := "[]"
	if len(req.OCRText) > 0 {
		joined := make([]string, 0, len(req.OCRText))
		for _, v := range req.OCRText {
			joined = append(joined, strings.ReplaceAll(v, "\"", "'"))
		}
		ocrTextJSON = "[\"" + strings.Join(joined, "\",\"") + "\"]"
	}

	feedback := &FeedbackSample{}
	err := db.conn.QueryRow(`
		INSERT INTO feedback_samples (
			part_number_predicted, part_number_final,
			serial_number_predicted, serial_number_final,
			manufacturer_predicted, manufacturer_final,
			category_predicted, category_final,
			correction_applied, confidence,
			image_data, ocr_text, meta_json
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id,
			COALESCE(part_number_predicted, ''),
			COALESCE(part_number_final, ''),
			COALESCE(serial_number_predicted, ''),
			COALESCE(serial_number_final, ''),
			COALESCE(manufacturer_predicted, ''),
			COALESCE(manufacturer_final, ''),
			COALESCE(category_predicted, ''),
			COALESCE(category_final, ''),
			correction_applied,
			COALESCE(confidence, 0),
			COALESCE(image_data, ''),
			COALESCE(meta_json, ''),
			created_at
	`,
		req.PartNumberPredicted,
		req.PartNumberFinal,
		req.SerialNumberPredicted,
		req.SerialNumberFinal,
		req.ManufacturerPredicted,
		req.ManufacturerFinal,
		req.CategoryPredicted,
		req.CategoryFinal,
		req.CorrectionApplied,
		req.Confidence,
		req.ImageData,
		ocrTextJSON,
		req.MetaJSON,
	).Scan(
		&feedback.ID,
		&feedback.PartNumberPredicted,
		&feedback.PartNumberFinal,
		&feedback.SerialNumberPredicted,
		&feedback.SerialNumberFinal,
		&feedback.ManufacturerPredicted,
		&feedback.ManufacturerFinal,
		&feedback.CategoryPredicted,
		&feedback.CategoryFinal,
		&feedback.CorrectionApplied,
		&feedback.Confidence,
		&feedback.ImageData,
		&feedback.MetaJSON,
		&feedback.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	feedback.OCRText = req.OCRText
	return feedback, nil
}

func (db *Database) ListFeedbackSamples(limit int, correctionsOnly bool) ([]FeedbackSample, int, error) {
	if limit <= 0 {
		limit = 100
	}

	baseWhere := ""
	args := []any{}
	if correctionsOnly {
		baseWhere = " WHERE correction_applied = TRUE"
	}

	var total int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM feedback_samples" + baseWhere).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id,
			COALESCE(part_number_predicted, ''),
			COALESCE(part_number_final, ''),
			COALESCE(serial_number_predicted, ''),
			COALESCE(serial_number_final, ''),
			COALESCE(manufacturer_predicted, ''),
			COALESCE(manufacturer_final, ''),
			COALESCE(category_predicted, ''),
			COALESCE(category_final, ''),
			correction_applied,
			COALESCE(confidence, 0),
			COALESCE(image_data, ''),
			COALESCE(ocr_text, '[]'),
			COALESCE(meta_json, ''),
			created_at
		FROM feedback_samples` + baseWhere + `
		ORDER BY created_at DESC
		LIMIT $1
	`
	args = append(args, limit)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []FeedbackSample{}
	for rows.Next() {
		item := FeedbackSample{}
		var ocrTextRaw string
		err = rows.Scan(
			&item.ID,
			&item.PartNumberPredicted,
			&item.PartNumberFinal,
			&item.SerialNumberPredicted,
			&item.SerialNumberFinal,
			&item.ManufacturerPredicted,
			&item.ManufacturerFinal,
			&item.CategoryPredicted,
			&item.CategoryFinal,
			&item.CorrectionApplied,
			&item.Confidence,
			&item.ImageData,
			&ocrTextRaw,
			&item.MetaJSON,
			&item.CreatedAt,
		)
		if err != nil {
			continue
		}

		if ocrTextRaw != "" && ocrTextRaw != "[]" {
			_ = json.Unmarshal([]byte(ocrTextRaw), &item.OCRText)
		}
		items = append(items, item)
	}

	return items, total, rows.Err()
}
