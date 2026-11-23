package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB inicializa la base de datos y crea las tablas si no existen
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario TEXT UNIQUE,
		correo TEXT UNIQUE,
		password_hash TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creando la tabla usuarios:", err)
	}

	createFavsTableSQL := `CREATE TABLE IF NOT EXISTS favoritos (
		usuario TEXT,
		clase_id TEXT,
		PRIMARY KEY (usuario, clase_id)
	);`
	_, err = db.Exec(createFavsTableSQL)
	if err != nil {
		log.Fatal("Error creando la tabla favoritos:", err)
	}

	createBossTableSQL := `CREATE TABLE IF NOT EXISTS jefes_derrotados (
		usuario TEXT,
		jefe_id TEXT,
		PRIMARY KEY (usuario, jefe_id)
	);`
	_, err = db.Exec(createBossTableSQL)
	if err != nil {
		log.Fatal("Error creando la tabla jefes_derrotados:", err)
	}

	createItemsTableSQL := `CREATE TABLE IF NOT EXISTS objetos_obtenidos (
		usuario TEXT,
		objeto_id TEXT,
		PRIMARY KEY (usuario, objeto_id)
	);`
	_, err = db.Exec(createItemsTableSQL)
	if err != nil {
		log.Fatal("Error creando la tabla objetos_obtenidos:", err)
	}

	createNotasTableSQL := `CREATE TABLE IF NOT EXISTS notas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario TEXT,
    titulo TEXT,
    contenido TEXT
	)`
	_, err = db.Exec(createNotasTableSQL)
	if err != nil {
		log.Fatal("Error creando la tabla notas:", err)
	}

	return db, nil
}
