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

func DBConnect() error {

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

var queries = map[string]string{
	"purge_payments": `TRUNCATE TABLE payments;`,
	"insert_payments": `
		INSERT INTO payments (correlation_id, amount, processor)
		VALUES ($1, $2, $3)
        RETURNING id ;`,
}

func ExecuteQuery(queryName string) error {

	query, ok := queries[queryName]
	if !ok {
		return fmt.Errorf("query %s not found", queryName)
	}

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error executing query %s: %w", queryName, err)
	}

	return nil

}

func PurgeTable() error {
	return ExecuteQuery("purge_payments")
}

func InsertPayment(correlation_id string, amount float64, processor string) error {

	query, ok := queries["insert_payments"]
	if !ok {
		return fmt.Errorf("query insert_payments not found")
	}

	_, err := DB.Exec(query, correlation_id, amount, processor)
	if err != nil {
		return fmt.Errorf("error executing query insert_payments: %w", err)
	}

	return nil

}

func GetPaymentSummary(fromTime, toTime time.Time) (map[string]map[string]interface{}, error) {

	query := `
		WITH daily_summary AS (
			SELECT 
				processor,
				DATE_TRUNC('day', created_at) as day,
				COUNT(*) as total_requests,
				SUM(amount)::float as total_amount
			FROM payments
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY processor, DATE_TRUNC('day', created_at)
		)
		SELECT 
			processor,
			COALESCE(SUM(total_requests), 0) as total_requests,
			COALESCE(SUM(total_amount), 0.0) as total_amount
		FROM daily_summary
		GROUP BY processor;
	`

	rows, err := DB.Query(query, fromTime, toTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]map[string]interface{}{
		"default": {
			"totalRequests": 0,
			"totalAmount":   0.0,
		},
		"fallback": {
			"totalRequests": 0,
			"totalAmount":   0.0,
		},
	}

	for rows.Next() {

		var processor string
		var totalRequests int
		var totalAmount float64

		err := rows.Scan(&processor, &totalRequests, &totalAmount)
		if err != nil {
			return nil, err
		}

		result[processor] = map[string]interface{}{
			"totalRequests": totalRequests,
			"totalAmount":   totalAmount,
		}
	}

	return result, nil

}
