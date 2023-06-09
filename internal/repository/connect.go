package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func ConnectToDB() *sql.DB {
	db, err := sql.Open("sqlite", "./repo.db")
	if err != nil {
		log.Fatal("failed to open database", err)
	}
	scheme := `
	CREATE TABLE IF NOT EXISTS statistics (
		endpoint STRING NOT NULL PRIMARY KEY,
		counter INTEGER
	) 
	`
	_, err = db.Exec(scheme)
	if err != nil {
		log.Fatal("failed to create table", err)
	}
	return db
}
