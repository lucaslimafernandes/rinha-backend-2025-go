package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lucaslimafernandes/rinha-backend-2025-go/models"
)

type Payment struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func PaymentSend(processor string, payment Payment) error {

	// PaymentWorker()

	default_uri := os.Getenv("PAYMENT_PROCESSOR_DEFAULT_URL")
	fallback_uri := os.Getenv("PAYMENT_PROCESSOR_FALLBACK_URL")

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	send := func(uri string, payment Payment) error {
		resp, err := client.Get(fmt.Sprintf("%v/payments", default_uri))
		if err != nil {
			// log.Printf("error: %v (%v, %v)\n", err, uri, payment)
			return fmt.Errorf("error: %v (%v, %v)", err, uri, payment)
		}
		defer resp.Body.Close()
		return nil
	}

	if processor == "default" {
		send(default_uri, payment)
	} else {
		send(fallback_uri, payment)
	}

	pp := models.Payment{
		Correlation_id: payment.CorrelationId,
		Amount:         payment.Amount,
		Processor:      processor,
	}

	paymentChan <- pp

	// err := models.InsertPayment(payment.CorrelationId, payment.Amount, processor)
	// if err != nil {
	// 	return fmt.Errorf("error inserting db: %v", err)
	// }

	return nil

}
