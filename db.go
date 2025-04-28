package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func initDB() (*sql.DB, error) {
	// Use in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping to ensure connection is established
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established (in-memory).")
	return db, nil
}

func buildSchema(db *sql.DB) error {
	sqlPath := filepath.Join("SQL", "build_tables.sql")
	log.Printf("Reading schema from: %s\n", sqlPath)

	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", sqlPath, err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}
	log.Println("Database schema built successfully.")
	return nil
}

func loadData(db *sql.DB) error {
	fileNames := []string{
		"Dzial", // Make sure order respects foreign key constraints if any
		"Grupa",
		"Klasa",
		"Kategoria",
		"PodKategoria",
		"Uprawa",
	}
	csvDir := "CSV"

	for _, name := range fileNames {
		csvPath := filepath.Join(csvDir, fmt.Sprintf("%s.csv", name))
		log.Printf("Loading data for table '%s' from: %s\n", name, csvPath)

		file, err := os.Open(csvPath)
		if err != nil {
			// Allow skipping if file not found for tables like 'Dzial' if not essential for main view
			log.Printf("Warning: Could not open CSV file %s: %v. Skipping.", csvPath, err)
			continue // Or return error if all files are mandatory
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = ';' // Set delimiter
		reader.LazyQuotes = true // Handle quotes flexibly

		// Read header
		headers, err := reader.Read()
		if err != nil {
			return fmt.Errorf("failed to read header from %s: %w", csvPath, err)
		}

		// Prepare insert statement
		placeholders := make([]string, len(headers))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			name, strings.Join(headers, ","), strings.Join(placeholders, ","))

		// Use transaction for bulk insert
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", name, err)
		}

		stmt, err := tx.Prepare(sqlQuery)
		if err != nil {
			tx.Rollback() // Rollback on prepare error
			return fmt.Errorf("failed to prepare statement for %s: %w", name, err)
		}
		defer stmt.Close()

		rowCount := 0
		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break // End of file
				}
				tx.Rollback() // Rollback on read error
				return fmt.Errorf("error reading data row from %s: %w", csvPath, err)
			}

			// Convert record (slice of strings) to slice of interface{}
			args := make([]interface{}, len(record))
			for i, v := range record {
				args[i] = v
			}

			_, err = stmt.Exec(args...)
			if err != nil {
				tx.Rollback() // Rollback on exec error
				return fmt.Errorf("failed to execute insert for %s, row %v: %w", name, record, err)
			}
			rowCount++
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction for %s: %w", name, err)
		}
		log.Printf("Successfully loaded %d rows into table '%s'.\n", rowCount, name)
	}

	log.Println("All data loaded successfully.")
	return nil
}