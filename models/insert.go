package models

import "fmt"

type payment struct {
	correlation_id string
	amount         float64
	processor      string
}

func BulkInsert(payments []payment) error {

	if len(payments) == 0 {
		return nil
	}

	query := `INSERT INTO payments (correlation_id, amount, processor) VALUES`
	values := []interface{}{}
	placeholders := ""

	for i, p := range payments {

		idx := i * 3
		placeholders += fmt.Sprintf("($%d, $%d, $%d)", idx+1, idx+2, idx+3)

		if i < len(payments)-1 {
			placeholders += ", "
		}
		values = append(values, p.correlation_id, p.amount, p.processor)

	}

	query += placeholders + ";"
	_, err := DB.Exec(query, values...)

	if err != nil {
		return err
	}

	return nil
}
