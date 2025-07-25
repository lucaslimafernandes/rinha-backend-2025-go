package utils

import (
	"api-rinha/models"
	"fmt"
)

type Payment struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func PaymentSend(payment Payment) error {

	err := models.SendKafka(payment.CorrelationId, payment.Amount)
	if err != nil {
		return fmt.Errorf("error inserting db: %v", err)
	}

	return nil

}
