package main

import "database/sql"

// Using pointers for fields that might be conceptually nullable or text fields
// SQLite often treats bools as integers (0/1), so use bool and handle scanning appropriately.
// Match field names to template usage later (or use JSON tags if planning APIs).

type Uprawa struct {
	IdUprawa                 string            `db:"IdUprawa"` // Assuming PK is integer
	IdPodKategoria           string            `db:"IdPodKategoria"`
	NazwaUprawa              sql.NullString `db:"NazwaUprawa"`
	NazwaLacinskaUprawa      sql.NullString `db:"NazwaLacinskaUprawa"`
	NazwaSynonimyUprawa      sql.NullString `db:"NazwaSynonimyUprawa"`
	OpisUprawa               sql.NullString `db:"OpisUprawa"`
	UwagaUprawa              sql.NullString `db:"UwagaUprawa"`
	ProduktRolny             bool           `db:"ProduktRolny"`
	UprawaMiododajna         bool           `db:"UprawaMiododajna"`
	UprawaEkologiczna        bool           `db:"UprawaEkologiczna"`
	UprawaEnergetyczna       bool           `db:"UprawaEnergetyczna"`
	UprawaOgrodnicza         bool           `db:"UprawaOgrodnicza"`
	DostawyBezposrednie      bool           `db:"DostawyBezposrednie"`
	RolniczyHandelDetaliczny bool           `db:"RolniczyHandelDetaliczny"`
	DzialSpecjalny           bool           `db:"DzialSpecjalny"`
	OkrywaZimowa             bool           `db:"OkrywaZimowa"`
	Warzywo                  bool           `db:"Warzywo"`
	WarzywoOwocKwiatZiolo    bool           `db:"WarzywoOwocKwiatZiolo"`
}

type PodKategoria struct {
	IdPodKategoria    string
	NazwaPodKategoria string
	Uprawy            []Uprawa
}

type Kategoria struct {
	IdKategoria    string
	NazwaKategoria string
	PodKategorie   []*PodKategoria // Pointer to allow modifications
}

type Klasa struct {
	IdKlasa    string
	NazwaKlasa string
	Kategorie  []*Kategoria // Pointer to allow modifications
}

type Grupa struct {
	IdGrupa    string
	NazwaGrupa string
	Klasy      []*Klasa // Pointer to allow modifications
}

// Structure to pass data to the main template
type IndexData struct {
	Headers []string
	Grupy   []*Grupa // Use pointers for consistency with map storage
}

// Python WTForms equivalent isn't directly available. We'll just pass headers for display.
// If form *submission* was needed, we'd parse request form values.
var uprawaHeaders = []string{
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