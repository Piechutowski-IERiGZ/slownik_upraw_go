package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"slownik_upraw/ui"
	"slownik_upraw/models"	

	_ "github.com/mattn/go-sqlite3"
)

// GLOBALS
var db *sql.DB
var uprawaHeaders = []string {
	"Id Uprawy",
	"Id Podkategorii",
	"Nazwa Uprawy",
	"Nazwa Łacińska Uprawy",
	"Synonimy Nazwy Uprawy",
	"Opis Uprawy",
	"Uwagi do Uprawy",
	"Produkt Rolny",
	"Uprawa Miododajna",
	"Uprawa Ekologiczna",
	"Uprawa Energetyczna",
	"Uprawa Ogrodnicza",
	"Dostawy Bezpośrednie",
	"Rolniczy Handel Detaliczny",
	"Dział Specjalny Produkcji Rolnej",
	"Okrywa Zimowa",
	"Warzywo",
	"Warzywo / Owoc / Kwiat / Zioło",
}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	query := `
    SELECT
        g.IdGrupa, g.NazwaGrupa,
        kla.IdKlasa, kla.NazwaKlasa,
        k.IdKategoria, k.NazwaKategoria,
        pk.IdPodKategoria, pk.NazwaPodKategoria,
        u.IdUprawa, u.IdPodKategoria, u.NazwaUprawa, u.NazwaLacinskaUprawa, u.NazwaSynonimyUprawa,
        u.OpisUprawa, u.UwagaUprawa, u.ProduktRolny, u.UprawaMiododajna, u.UprawaEkologiczna,
        u.UprawaEnergetyczna, u.UprawaOgrodnicza, u.DostawyBezposrednie, u.RolniczyHandelDetaliczny,
        u.DzialSpecjalny, u.OkrywaZimowa, u.Warzywo, u.WarzywoOwocKwiatZiolo
    FROM
        Uprawa AS u
    JOIN
        PodKategoria AS pk ON u.IdPodKategoria = pk.IdPodKategoria
    JOIN
        Kategoria AS k ON pk.IdKategoria = k.IdKategoria
    JOIN
        Klasa AS kla ON k.IdKlasa = kla.IdKlasa
    JOIN
        Grupa AS g ON kla.IdGrupa = g.IdGrupa
    ORDER BY
        g.NazwaGrupa, kla.NazwaKlasa, k.NazwaKategoria, pk.NazwaPodKategoria, u.NazwaUprawa
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	grupaMap := make(map[string]*models.Grupa)
	klasaMap := make(map[string]*models.Klasa)
	kategoriaMap := make(map[string]*models.Kategoria)
	podkategoriaMap := make(map[string]*models.PodKategoria)

	for rows.Next() {
		var gID, klaID, kID, pkID, uID, uPKID,
			gName, klaName, kName, pkName string
		var u models.Uprawa
		var sqlProduktRolny, sqlUprawaMiododajna, sqlUprawaEkologiczna, sqlUprawaEnergetyczna, sqlUprawaOgrodnicza,
			sqlDostawyBezposrednie, sqlRolniczyHandelDetaliczny, sqlDzialSpecjalny, sqlOkrywaZimowa, sqlWarzywo,
			sqlWarzywoOwocKwiatZiolo sql.NullString

		err := rows.Scan(
			&gID, &gName, &klaID, &klaName, &kID, &kName, &pkID, &pkName,
			&uID, &uPKID, &u.NazwaUprawa, &u.NazwaLacinskaUprawa, &u.NazwaSynonimyUprawa,
			&u.OpisUprawa, &u.UwagaUprawa,
			&sqlProduktRolny,
			&sqlUprawaMiododajna,
			&sqlUprawaEkologiczna,
			&sqlUprawaEnergetyczna,
			&sqlUprawaOgrodnicza,
			&sqlDostawyBezposrednie,
			&sqlRolniczyHandelDetaliczny,
			&sqlDzialSpecjalny,
			&sqlOkrywaZimowa,
			&sqlWarzywo,
			&sqlWarzywoOwocKwiatZiolo,
		)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		u.IdUprawa = uID
		u.IdPodKategoria = uPKID
		u.ProduktRolny = sqlProduktRolny.Valid && strings.EqualFold(sqlProduktRolny.String, "true")
		u.UprawaMiododajna = sqlUprawaMiododajna.Valid && strings.EqualFold(sqlUprawaMiododajna.String, "true")
		u.UprawaEkologiczna = sqlUprawaEkologiczna.Valid && strings.EqualFold(sqlUprawaEkologiczna.String, "true")
		u.UprawaEnergetyczna = sqlUprawaEnergetyczna.Valid && strings.EqualFold(sqlUprawaEnergetyczna.String, "true")
		u.UprawaOgrodnicza = sqlUprawaOgrodnicza.Valid && strings.EqualFold(sqlUprawaOgrodnicza.String, "true")
		u.DostawyBezposrednie = sqlDostawyBezposrednie.Valid && strings.EqualFold(sqlDostawyBezposrednie.String, "true")
		u.RolniczyHandelDetaliczny = sqlRolniczyHandelDetaliczny.Valid && strings.EqualFold(sqlRolniczyHandelDetaliczny.String, "true")
		u.DzialSpecjalny = sqlDzialSpecjalny.Valid && strings.EqualFold(sqlDzialSpecjalny.String, "true")
		u.OkrywaZimowa = sqlOkrywaZimowa.Valid && strings.EqualFold(sqlOkrywaZimowa.String, "true")
		u.Warzywo = sqlWarzywo.Valid && strings.EqualFold(sqlWarzywo.String, "true")
		u.WarzywoOwocKwiatZiolo = sqlWarzywoOwocKwiatZiolo.Valid && strings.EqualFold(sqlWarzywoOwocKwiatZiolo.String, "true")

		grupa, ok := grupaMap[gID]
		if !ok {
			grupa = &models.Grupa{IdGrupa: gID, NazwaGrupa: gName}
			grupaMap[gID] = grupa
		}

		klasa, ok := klasaMap[klaID]
		if !ok {
			klasa = &models.Klasa{IdKlasa: klaID, NazwaKlasa: klaName}
			klasaMap[klaID] = klasa
			grupa.Klasy = append(grupa.Klasy, klasa)
		}

		kategoria, ok := kategoriaMap[kID]
		if !ok {
			kategoria = &models.Kategoria{IdKategoria: kID, NazwaKategoria: kName}
			kategoriaMap[kID] = kategoria
			klasa.Kategorie = append(klasa.Kategorie, kategoria)
		}

		podkategoria, ok := podkategoriaMap[pkID]
		if !ok {
			podkategoria = &models.PodKategoria{IdPodKategoria: pkID, NazwaPodKategoria: pkName}
			podkategoriaMap[pkID] = podkategoria
			kategoria.PodKategorie = append(kategoria.PodKategorie, podkategoria)
		}

		podkategoria.Uprawy = append(podkategoria.Uprawy, u)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	grupySlice := make([]*models.Grupa, 0, len(grupaMap))
	for _, g := range grupaMap {
		grupySlice = append(grupySlice, g)
	}
	sort.Slice(grupySlice, func(i, j int) bool {
		return grupySlice[i].NazwaGrupa < grupySlice[j].NazwaGrupa
	})

	data := models.IndexData{
		Headers: uprawaHeaders,
		Grupy:   grupySlice,
	}

	ui.SlownikPage(data).Render(r.Context(), w)
	// err = templates.ExecuteTemplate(w, "slownik.gohtml", data)
	// if err != nil {
	// 	log.Printf("Error executing template slownik.gohtml: %v\n", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatalf("Failed to poing DB: %v", err)
	}

	log.Println("Database connection established (in-memory).")
}

func buildSchema() {
	sqlPath := filepath.Join("SQL", "build_tables.sql")
	sqlBytes, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatalf("Failed to read schema file %v: %v", sqlPath, err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatal("failed to execute schema SQL: %w", err)
	}
	log.Println("Database schema built successfully.")
}

func loadData() {
	fileNames := []string{
		"Dzial",
		"Grupa",
		"Klasa",
		"Kategoria",
		"PodKategoria",
		"Uprawa",

		// "ZmianowanieWarzywo",
		// "WykorzystanieProdukt",
		// "TypUprawa",
		// "RodzajUprawa",

		// "Uprawa_ZmianowanieWarzywo",
		// "Uprawa_WykorzystanieProdukt",
		// "Uprawa_TypUprawa",
		// "Uprawa_RodzajUprawa",
	}
	csvDir := "CSV"

	for _, name := range fileNames {
		csvPath := filepath.Join(csvDir, fmt.Sprintf("%v.csv", name))
		file, err := os.Open(csvPath)
		if err != nil {
			log.Fatalf("Warning: Could not open CSV file %s: %v", csvPath, err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = ';'
		reader.LazyQuotes = true

		headers, err := reader.Read()
		if err != nil {
			log.Fatalf("failed to read header from %s: %v", csvPath, err)
		}

		placeholders := make([]string, len(headers))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		sqlQuery := fmt.Sprintf(
			"INSERT INTO %v (%v) VALUES (%v)",
			name,
			strings.Join(headers, ","),
			strings.Join(placeholders, ","),
		)

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction for %v: %v", name, err)
		}

		stmt, err := tx.Prepare(sqlQuery)
		if err != nil {
			log.Fatalf("failed to prepare statement for %v: %v", name, err)
		}
		defer stmt.Close()

		rowCount := 0
		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Fatalf("Error reading data row from %v: %v", csvPath, err)
			}

			args := make([]interface{}, len(record))
			for i, v := range record {
				args[i] = v
			}

			_, err = stmt.Exec(args...)
			if err != nil {
				log.Fatalf("Failed to execute insert for %v, row %v: %v", name, record, err)
			}
			rowCount++
		}

		err = tx.Commit()
		if err != nil {
			log.Fatalf("failed to commit transaction for %v: %v", name, err)
		}
		log.Printf("Successfully loaded %v rows into table '%v'.\n", rowCount, name)
	}

	log.Println("All data loaded successfully.")
}


func main() {
	initDB()
	buildSchema()
	loadData()
	defer db.Close()
	// initTemplates()

	// Setup HTTP Server
	mux := http.NewServeMux()

	// Static files
	staticDir := http.Dir("Static")
	staticServer := http.FileServer(staticDir)
	mux.Handle("/Static/", http.StripPrefix("/Static/", staticServer))
	log.Println("Serving Static files from /Static/ route.")

	// Application Routes
	mux.HandleFunc("/", indexHandler)

	// Start Server
	port := "8080"
	log.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
