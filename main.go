package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

var (
	templates *template.Template
	db        *sql.DB
)

func main() {
	// Initialize Database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close() // Ensure DB connection is closed on exit

	// Build Schema
	if err := buildSchema(db); err != nil {
		log.Fatalf("Failed to build database schema: %v", err)
	}

	// Load Data
	if err := loadData(db); err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Parse Templates
	// Use ParseGlob to load all templates in the directory
	// Ensure template names used in ExecuteTemplate match the filenames (e.g., "slownik.gohtml")
	templatePattern := filepath.Join("Templates", "*.gohtml")
	templates = template.Must(template.ParseGlob(templatePattern))
	log.Printf("Templates parsed from %s\n", templatePattern)

	// Setup HTTP Server
	mux := http.NewServeMux()

	// Static files
	staticDir := http.Dir("Static")
	staticServer := http.FileServer(staticDir)
	mux.Handle("/Static/", http.StripPrefix("/Static/", staticServer))
	log.Println("Serving Static files from /Static/ route.")

	// Application Routes
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/download-form", downloadFormHandler) // Note: Renamed to avoid clash
	mux.HandleFunc("/jak-to-dziala", jakToDzialaHandler)
	mux.HandleFunc("/api", apiHandler)

	// Start Server
	port := "8080"
	log.Printf("Starting server on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// --- Handlers ---

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

	// --- CHANGE MAP KEY TYPE ---
	// Use maps with STRING keys for hierarchy, pointers allow modification
	grupaMap := make(map[string]*Grupa)
	klasaMap := make(map[string]*Klasa)
	kategoriaMap := make(map[string]*Kategoria)
	podkategoriaMap := make(map[string]*PodKategoria)
	// --- END CHANGE ---

	for rows.Next() {
		// --- USE STRING FOR IDS ---
		var gID, klaID, kID, pkID, uID, uPKID string
		var gName, klaName, kName, pkName string
		// --- END CHANGE ---
		var u Uprawa // Assume Uprawa struct fields for IDs are now also string
		// --- **CRITICAL FIX: USE sql.NullString FOR BOOLEANS** ---
		var sqlProduktRolny, sqlUprawaMiododajna, sqlUprawaEkologiczna, sqlUprawaEnergetyczna, sqlUprawaOgrodnicza, sqlDostawyBezposrednie, sqlRolniczyHandelDetaliczny, sqlDzialSpecjalny, sqlOkrywaZimowa, sqlWarzywo, sqlWarzywoOwocKwiatZiolo sql.NullString
		
		err := rows.Scan(
			&gID, &gName, &klaID, &klaName, &kID, &kName, &pkID, &pkName,
			&uID, &uPKID, &u.NazwaUprawa, &u.NazwaLacinskaUprawa, &u.NazwaSynonimyUprawa,
			&u.OpisUprawa, &u.UwagaUprawa,
			&sqlProduktRolny, // Scan into sql.NullString
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

        // Assign STRING IDs to Uprawa struct
		u.IdUprawa = uID
		u.IdPodKategoria = uPKID

        // Convert NullStrings to bools (this logic is now correct as variables are NullString)
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

		// Build Hierarchy using STRING IDs as map keys
		grupa, ok := grupaMap[gID]
		if !ok {
			// Assign STRING ID when creating struct
			grupa = &Grupa{IdGrupa: gID, NazwaGrupa: gName}
			grupaMap[gID] = grupa
		}

		klasa, ok := klasaMap[klaID]
		if !ok {
			// Assign STRING ID when creating struct
			klasa = &Klasa{IdKlasa: klaID, NazwaKlasa: klaName}
			klasaMap[klaID] = klasa
			grupa.Klasy = append(grupa.Klasy, klasa)
		}

		kategoria, ok := kategoriaMap[kID]
		if !ok {
			// Assign STRING ID when creating struct
			kategoria = &Kategoria{IdKategoria: kID, NazwaKategoria: kName}
			kategoriaMap[kID] = kategoria
			klasa.Kategorie = append(klasa.Kategorie, kategoria)
		}

		podkategoria, ok := podkategoriaMap[pkID]
		if !ok {
			// Assign STRING ID when creating struct
			podkategoria = &PodKategoria{IdPodKategoria: pkID, NazwaPodKategoria: pkName}
			podkategoriaMap[pkID] = podkategoria
			kategoria.PodKategorie = append(kategoria.PodKategorie, podkategoria)
		}

		// Add the Uprawa to the current PodKategoria
		podkategoria.Uprawy = append(podkategoria.Uprawy, u)
	} // End for rows.Next()

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating rows: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    // Extract sorted Grupy from map (sorting by NazwaGrupa remains unchanged)
    grupySlice := make([]*Grupa, 0, len(grupaMap))
    for _, g := range grupaMap {
        grupySlice = append(grupySlice, g)
    }
    sort.Slice(grupySlice, func(i, j int) bool {
		return grupySlice[i].NazwaGrupa < grupySlice[j].NazwaGrupa
	})


	data := IndexData{
		Headers: uprawaHeaders, // Ensure uprawaHeaders is defined globally or passed appropriately
		Grupy:   grupySlice,
	}

	// Execute the main template
	err = templates.ExecuteTemplate(w, "slownik.gohtml", data)
	if err != nil {
		log.Printf("Error executing template slownik.gohtml: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
} 

func testHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "test.gohtml", nil) // Pass nil data
	if err != nil {
		log.Printf("Error executing template test.gohtml: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func jakToDzialaHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "jak_to_dziala.gohtml", nil)
	if err != nil {
		log.Printf("Error executing template jak_to_dziala.gohtml: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "api.gohtml", nil)
	if err != nil {
		log.Printf("Error executing template api.gohtml: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func downloadFormHandler(w http.ResponseWriter, r *http.Request) {
	// Same hardcoded data as the Python example
	data := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "30", "New York"},
		{"Bob", "25", "London"},
		{"Charlie", "35", "Paris"},
	}

	var buf bytes.Buffer
	// Explicitly create writer with buffer is better for capturing output
	csvWriter := csv.NewWriter(&buf)

	// Write header separately if needed, or use WriteAll
	// if err := csvWriter.Write(data[0]); err != nil { ... }
	// for _, record := range data[1:] { if err := csvWriter.Write(record); ... }

	// WriteAll is simpler for this case
	if err := csvWriter.WriteAll(data); err != nil {
		log.Printf("Error writing CSV data: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	csvWriter.Flush() // Don't forget to flush!
	if err := csvWriter.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=download.csv")

	// Write buffer to response
	w.Write(buf.Bytes())
}