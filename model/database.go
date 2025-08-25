package model

import (
	"database/sql"
	"log"
)

func SetupDatabase(db *sql.DB) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS ltp (
		pair TEXT PRIMARY KEY,
		amount REAL
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert mock data
	mockData := map[string]float64{
		"BTC/USD": 52000.12,
		"BTC/CHF": 49000.12,
		"BTC/EUR": 50000.12,
	}
	for pair, amount := range mockData {
		_, err := db.Exec(`INSERT OR REPLACE INTO ltp (pair, amount) VALUES (?, ?)`, pair, amount)
		if err != nil {
			log.Fatalf("Failed to insert mock data: %v", err)
		}
	}
}

func GetLTP(db *sql.DB, pair string) (float64, error) {
	var amount float64
	row := db.QueryRow(`SELECT amount FROM ltp WHERE pair = ?`, pair)
	if err := row.Scan(&amount); err != nil {
		return 0, err
	}
	return amount, nil
}
