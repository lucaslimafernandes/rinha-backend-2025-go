package models

// type Payment struct {
// 	Correlation_id string
// 	Amount         float64
// 	Processor      string
// }

// // func BulkInsert(payments []Payment) error {

// // 	if len(payments) == 0 {
// // 		return nil
// // 	}

// // 	query := `INSERT INTO payments (correlation_id, amount, processor) VALUES`
// // 	values := []interface{}{}
// // 	placeholders := ""

// // 	for i, p := range payments {

// // 		idx := i * 3
// // 		placeholders += fmt.Sprintf("($%d, $%d, $%d)", idx+1, idx+2, idx+3)

// // 		if i < len(payments)-1 {
// // 			placeholders += ", "
// // 		}
// // 		values = append(values, p.Correlation_id, p.Amount, p.Processor)

// // 	}

// // 	query += placeholders + ";"
// // 	_, err := DB.Exec(query, values...)

// // 	return err
// // }

// func BulkInsert(payments []Payment) error {
// 	if len(payments) == 0 {
// 		return nil
// 	}

// 	rows := make([][]interface{}, len(payments))
// 	for i, p := range payments {
// 		rows[i] = []interface{}{p.Correlation_id, p.Amount, p.Processor}
// 	}

// 	_, err := DB.CopyFrom(
// 		context.Background(),
// 		pgx.Identifier{"payments"},
// 		[]string{"correlation_id", "amount", "processor"},
// 		pgx.CopyFromRows(rows),
// 	)

// 	return err
// }
