package utils

import (
	"log"
	"sync"
	"time"

	"github.com/lucaslimafernandes/rinha-backend-2025-go/models"
)

const (
	BATCH_SIZE  = 100
	TIME_TICKER = 2000 * time.Millisecond
)

var (
	paymentChan = make(chan models.Payment, 500)
	once        sync.Once
)

func PaymentWorker() {

	log.Printf("BATCH_SIZE: %v\nTIME_TICKER: %v\n", BATCH_SIZE, TIME_TICKER)

	once.Do(func() {
		go func() {

			batch := make([]models.Payment, 0, BATCH_SIZE)
			ticker := time.NewTicker(TIME_TICKER)
			defer ticker.Stop()

			for {
				select {
				case p := <-paymentChan:
					batch = append(batch, p)
					if len(batch) >= BATCH_SIZE {
						go models.BulkInsert(batch)
						batch = make([]models.Payment, 0, BATCH_SIZE)
					}

				case <-ticker.C:
					if len(batch) > 0 {
						go models.BulkInsert(batch)
						batch = make([]models.Payment, 0, BATCH_SIZE)

					}
				}
			}
		}()

	})

}
