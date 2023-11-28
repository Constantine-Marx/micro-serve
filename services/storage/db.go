package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(user, password, host, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	log.Println("Connected to database")
	return db, nil
}
