package main

import (
	"database/sql"
	"log"
)

func bootstrapSQLite(conn *sql.DB) error {
	statements := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS catalog (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			part_number TEXT UNIQUE NOT NULL,
			serial_pattern TEXT,
			manufacturer TEXT,
			category TEXT,
			normalized_description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			catalog_id INTEGER NOT NULL REFERENCES catalog(id) ON DELETE CASCADE,
			serial_number TEXT,
			quantity INTEGER DEFAULT 1,
			location TEXT,
			status TEXT DEFAULT 'active',
			received_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS movements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			inventory_id INTEGER NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL,
			operation TEXT NOT NULL,
			reason TEXT,
			user_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_catalog_pn ON catalog(part_number);`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_status ON inventory(status);`,
	}

	for _, stmt := range statements {
		if _, err := conn.Exec(stmt); err != nil {
			log.Printf("Erro no bootstrap SQLite: %v", err)
			return err
		}
	}

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM catalog`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seed := []struct {
		pn, manufacturer, category, description string
	}{
		{"M393A4K40DB3-CWE", "Samsung", "memory", "DDR4 32GB 3200MHz RDIMM ECC"},
		{"HUH721212AL5200", "HGST", "disk", "HDD 12TB SAS 7200RPM 256MB Cache"},
		{"INTEL_SSDPE2MX450G7", "Intel", "ssd", "NVMe SSD 450GB 3D TLC PCI-E 3.0"},
		{"E82968-001", "HP", "network", "NIC 10Gb Ethernet Adapter Dual Port"},
		{"UCSC-PCIE-QSFP28", "Cisco", "network", "QSFP28 Module Adapter"},
	}

	for _, item := range seed {
		if _, err := conn.Exec(`
			INSERT INTO catalog (part_number, manufacturer, category, normalized_description)
			VALUES ($1, $2, $3, $4)
		`, item.pn, item.manufacturer, item.category, item.description); err != nil {
			return err
		}
	}

	return nil
}
