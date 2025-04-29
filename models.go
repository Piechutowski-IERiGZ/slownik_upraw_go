package main

import "database/sql"


type Uprawa struct {
	IdUprawa                 string         `db:"IdUprawa"`
	IdPodKategoria           string         `db:"IdPodKategoria"`
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
	PodKategorie   []*PodKategoria
}

type Klasa struct {
	IdKlasa    string
	NazwaKlasa string
	Kategorie  []*Kategoria
}

type Grupa struct {
	IdGrupa    string
	NazwaGrupa string
	Klasy      []*Klasa 
}

type IndexData struct {
	Headers []string
	Grupy   []*Grupa
}


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