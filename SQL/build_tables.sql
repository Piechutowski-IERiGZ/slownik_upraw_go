CREATE TABLE "Dzial" (
  "IdDzial" TEXT PRIMARY KEY NOT NULL,
  "NazwaDzial" TEXT,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT
);

CREATE TABLE "Grupa" (
  "IdGrupa" TEXT PRIMARY KEY NOT NULL,
  "IdDzial" TEXT,
  "NazwaGrupa" TEXT,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT,
  FOREIGN KEY ("IdDzial") REFERENCES "Dzial" ("IdDzial")
);

CREATE TABLE "Klasa" (
  "IdKlasa" TEXT PRIMARY KEY NOT NULL,
  "IdGrupa" TEXT,
  "NazwaKlasa" TEXT,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT,
  FOREIGN KEY ("IdGrupa") REFERENCES "Grupa" ("IdGrupa")
);

CREATE TABLE "Kategoria" (
  "IdKategoria" TEXT PRIMARY KEY NOT NULL,
  "IdKlasa" TEXT,
  "NazwaKategoria" TEXT,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT,
  FOREIGN KEY ("IdKlasa") REFERENCES "Klasa" ("IdKlasa")
);

CREATE TABLE "PodKategoria" (
  "IdPodKategoria" TEXT PRIMARY KEY NOT NULL,
  "IdKategoria" TEXT,
  "NazwaPodKategoria" TEXT,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT,
  FOREIGN KEY ("IdKategoria") REFERENCES "Kategoria" ("IdKategoria")
);

CREATE TABLE "Uprawa" (
  "IdUprawa" TEXT PRIMARY KEY NOT NULL,
  "IdPodKategoria" TEXT,
  "NazwaUprawa" TEXT,
  "NazwaLacinskaUprawa" TEXT,
  "NazwaSynonimyUprawa" TEXT,
  "OpisUprawa" TEXT,
  "UwagaUprawa" TEXT,
  "IdProdukt" TEXT,
  "ProduktRolny" BOOLEAN,
  "UprawaMiododajna" BOOLEAN,
  "UprawaEkologiczna" BOOLEAN,
  "UprawaEnergetyczna" BOOLEAN,
  "UprawaOgrodnicza" BOOLEAN,
  "DostawyBezposrednie" BOOLEAN,
  "RolniczyHandelDetaliczny" BOOLEAN,
  "DzialSpecjalny" BOOLEAN,
  "OkrywaZimowa" BOOLEAN,
  "Warzywo" BOOLEAN,
  "WarzywoOwocKwiatZiolo" BOOLEAN,
  "IdStatus" TEXT,
  "Status" TEXT,
  "DataDodania" TIMESTAMP,
  "DodanyPrzez" TEXT,
  "DataModyfikacji" TIMESTAMP,
  "ZmodyfikowanyPrzez" TEXT,
  "DataUsuniecia" TIMESTAMP,
  "UsunietyPrzez" TEXT,
  FOREIGN KEY ("IdPodKategoria") REFERENCES "PodKategoria" ("IdPodKategoria")
);
