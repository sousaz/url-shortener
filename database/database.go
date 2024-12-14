package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sousaz/urlshortener/utils"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./app.db")
	if err != nil {
		fmt.Println("Erro initializing database")
	}

	_, err = DB.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS urls (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		original TEXT UNIQUE NOT NULL,
		shortener TEXT NOT NULL UNIQUE
		)`,
	)

	if err != nil {
		fmt.Println("Error creating table")
	}
}

func AddUrl(u map[string]interface{}) (*string, error) {
	shortener := utils.Generate_shortener(8)
	_, err := DB.ExecContext(
		context.Background(),
		`
		INSERT OR IGNORE INTO urls (original, shortener)
		VALUES (?, ?);
		`, u["original"], shortener,
	)
	if err != nil {
		return nil, err
	}

	var savedShortener string
	row := DB.QueryRowContext(
		context.Background(),
		`
		SELECT shortener FROM urls WHERE original = (?);
		`, u["original"],
	)
	err = row.Scan(&savedShortener)
	if err != nil {
		return nil, err
	}

	return &savedShortener, nil
}

func GetUrl(id string) (*string, error) {
	var original string
	row := DB.QueryRowContext(
		context.Background(),
		`
		SELECT original FROM urls WHERE shortener = (?);
		`, id,
	)
	err := row.Scan(&original)
	if err != nil {
		return nil, err
	}
	return &original, nil
}
