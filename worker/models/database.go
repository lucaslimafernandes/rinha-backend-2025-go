package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var insert_payment = `
	INSERT INTO payments (correlation_id, amount, processor)
	VALUES ($1, $2, $3)
	RETURNING id ;`

func ConnectDB() error {

	var err error

	connStr := os.Getenv("PG_DSN")

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect: %v\n", err)
		return err
	}

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Second)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping: %v\n", err)
	}

	return nil

}

func InsertPayment(correlation_id string, amount float64, processor string) error {

	_, err := DB.Exec(insert_payment, correlation_id, amount, processor)
	if err != nil {
		return fmt.Errorf("error executing query insert_payments: %w", err)
	}

	return nil

}
