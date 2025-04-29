package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	_ "github.com/mattn/go-sqlite3" 
)

func initDB() (*sql.DB ) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)		
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatalf("Failed to poing DB: %v", err)
	}

	log.Println("Database connection established (in-memory).")
	return db
}

func buildSchema(db *sql.DB) {
	sqlPath := filepath.Join("SQL", "build_tables.sql")
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatal("Failed to read schema file %s: %w", sqlPath, err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("failed to execute schema SQL: %w", err)
	}
	log.Println("Database schema built successfully.")
}

func loadData(db *sql.DB) {
	fileNames := []string{
		"Dzial", 
		"Grupa",
		"Klasa",
		"Kategoria",
		"PodKategoria",
		"Uprawa",

		"ZmianowanieWarzywo",
		"WykorzystanieProdukt",
		"TypUprawa",
		"RodzajUprawa",
	
		"Uprawa_ZmianowanieWarzywo",
		"Uprawa_WykorzystanieProdukt",
		"Uprawa_TypUprawa",
		"Uprawa_RodzajUprawa",	
	}
	csvDir := "CSV"

	for name := range fileNames {
		csvPath := filepath.Join(csvDir, fmt.Sprintf("%s.csv", name))
		file, err := os.Open(csvPath)
		if err != nil {
			log.Fatal("Warning: Could not open CSV file %s: %v", csvPath, err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.LazyQuotes = true 

		headers, err := reader.Read()
		if err != nil {
			log.Fatal("failed to read header from %s: %w", csvPath, err)
		}

		placeholders := make([]string, len(headers))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		sqlQuery := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			name, 
			strings.Join(headers, ","),
			strings.Join(placeholders, ","),
		)

		tx, err := db.Begin()
		if err != nil {
			log.Fatal("Failed to begin transaction for %s: %w", name, err)
		}

		stmt, err := tx.Prepare(sqlQuery)
		if err != nil {
			log.Fatal("failed to prepare statement for %s: %w", name, err)
		}
		defer stmt.Close()

		rowCount := 0
		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break 
				}
				log.Fatal("Error reading data row from %s: %w", csvPath, err)
			}

			args := make([]interface{}, len(record))
			for i, v := range record {
				args[i] = v
			}

			_, err = stmt.Exec(args...)
			if err != nil {
				log.Fatal("Failed to execute insert for %s, row %v: %w", name, record, err)
			}
			rowCount++
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal("failed to commit transaction for %s: %w", name, err)
		}
		log.Printf("Successfully loaded %d rows into table '%s'.\n", rowCount, name)
	}

	log.Println("All data loaded successfully.")
}