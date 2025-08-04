package models

// var DB *pgxpool.Pool

// func DBConnect() error {

// 	var err error

// 	connStr := os.Getenv("PG_DSN")
// 	cfg, err := pgxpool.ParseConfig(connStr)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse pool config: %v", err)
// 	}

// 	// pool configs
// 	cfg.MaxConns = 30
// 	cfg.MinConns = 5
// 	cfg.MaxConnLifetime = 5 * time.Minute

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	DB, err = pgxpool.NewWithConfig(ctx, cfg)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect: %v", err)
// 	}

// 	err = DB.Ping(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to ping: %v", err)
// 	}

// 	return nil

// }

// var queries = map[string]string{
// 	"purge_payments": `TRUNCATE TABLE payments;`,
// 	"insert_payments": `
// 		INSERT INTO payments (correlation_id, amount, processor)
// 		VALUES ($1, $2, $3)
//         RETURNING id ;`,
// }

// func ExecuteQuery(queryName string) error {

// 	query, ok := queries[queryName]
// 	if !ok {
// 		return fmt.Errorf("query %s not found", queryName)
// 	}

// 	_, err := DB.Exec(context.Background(), query)
// 	if err != nil {
// 		return fmt.Errorf("error executing query %s: %w", queryName, err)
// 	}

// 	return nil

// }

// func PurgeTable() error {
// 	return ExecuteQuery("purge_payments")
// }

// func InsertPayment(correlation_id string, amount float64, processor string) error {

// 	query, ok := queries["insert_payments"]
// 	if !ok {
// 		return fmt.Errorf("query insert_payments not found")
// 	}

// 	_, err := DB.Exec(context.Background(), query, correlation_id, amount, processor)
// 	if err != nil {
// 		return fmt.Errorf("error executing query insert_payments: %w", err)
// 	}

// 	return nil

// }

// func GetPaymentSummary(fromTime, toTime time.Time) (map[string]map[string]interface{}, error) {

// 	query := `
// 		SELECT
// 			processor,
// 			COUNT(*) AS total_requests,
// 			SUM(amount)::float AS total_amount
// 		FROM payments
// 		WHERE created_at BETWEEN $1 AND $2
// 		GROUP BY processor;
// 	`

// 	rows, err := DB.Query(context.Background(), query, fromTime, toTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	result := map[string]map[string]interface{}{
// 		"default": {
// 			"totalRequests": 0,
// 			"totalAmount":   0.0,
// 		},
// 		"fallback": {
// 			"totalRequests": 0,
// 			"totalAmount":   0.0,
// 		},
// 	}

// 	for rows.Next() {

// 		var processor string
// 		var totalRequests int
// 		var totalAmount float64

// 		err := rows.Scan(&processor, &totalRequests, &totalAmount)
// 		if err != nil {
// 			return nil, err
// 		}

// 		result[processor] = map[string]interface{}{
// 			"totalRequests": totalRequests,
// 			"totalAmount":   totalAmount,
// 		}
// 	}

// 	return result, nil

// }
