package service

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"worker-rinha/models"
)

func PaymentSend(processor string, payment models.Payment) error {

	default_uri := os.Getenv("PAYMENT_PROCESSOR_DEFAULT_URL")
	fallback_uri := os.Getenv("PAYMENT_PROCESSOR_FALLBACK_URL")

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	send := func(uri string, payment models.Payment) error {
		resp, err := client.Get(fmt.Sprintf("%v/payments", payment))
		if err != nil {
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

	err := models.InsertPayment(payment.CorrelationId, payment.Amount, processor)
	if err != nil {
		return fmt.Errorf("error inserting db: %v", err)
	}

	return nil

}
