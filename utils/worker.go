package utils

import (
	"log"
	"sync"

	"github.com/lucaslimafernandes/rinha-backend-2025-go/models"
)

const (
	BATCH_SIZE  = 100
	workerCount = 20
)

var (
	paymentChan = make(chan models.Payment, 500)
	// paymentChan = make(chan models.Payment, 10000)
	once sync.Once
)

// func paymentWorker() {

// 	once.Do(func() {
// 		go func() {

// 			batch := make([]models.Payment, 0, BATCH_SIZE)
// 			ticker := time.NewTicker(2 * time.Second)
// 			defer ticker.Stop()

// 			for {
// 				select {
// 				case p := <-paymentChan:
// 					batch = append(batch, p)
// 					if len(batch) >= BATCH_SIZE {
// 						go models.BulkInsert(batch)
// 						batch = make([]models.Payment, 0, BATCH_SIZE)
// 					}

// 				case <-ticker.C:
// 					if len(batch) > 0 {
// 						go models.BulkInsert(batch)
// 						batch = make([]models.Payment, 0, BATCH_SIZE)

// 					}
// 				}
// 			}
// 		}()

// 	})

// }

func paymentWorker() {

	once.Do(func() {
		log.Printf("Iniciando pool de %d workers...", workerCount)
		go func() {

			for i := 0; i < workerCount; i++ {
				go func(workerID int) {

					for p := range paymentChan {
						err := models.SimpleInsert(p)
						if err != nil {
							log.Printf("[Worker %d] erro no insert: %v", workerID, err)
						}
					}

				}(i)
			}

		}()

	})

}
