package service

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func A() {

	default_uri := os.Getenv("PAYMENT_PROCESSOR_DEFAULT_URL")
	fallback_uri := os.Getenv("PAYMENT_PROCESSOR_FALLBACK_URL")

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	send := func(uri string, payment Payment) error {
		resp, err := client.Get(fmt.Sprintf("%v/payments", payment))
		if err != nil {
			return fmt.Errorf("error: %v (%v, %v)", err, uri, payment)
		}
		defer resp.Body.Close()
		return nil
	}

}
