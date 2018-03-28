package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// DB is the database connection global variable
var DB *sql.DB

// Init initializes the DB connection
func Init() error {
	connConfig := `
		host=database
		user=coffee
		dbname=coffeeshop
		password=needcaffeine
		sslmode=disable
	`

	db, err := sql.Open("postgres", connConfig)
	if err != nil {
		return err
	}

	DB = db

	return nil
}

// Tx begins and returns the DB transaction
func Tx() (*sql.Tx, error) {
	return DB.Begin()
}
