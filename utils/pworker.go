package utils

// const (
// 	BATCH_SIZE     = 100
// 	TIME_TICKER    = 50 * time.Millisecond
// 	WORKER_COUNT   = 4
// 	PAYMENT_BUFFER = 5000
// 	BATCH_BUFFER   = 100
// )

// var (
// 	paymentChan = make(chan models.Payment, PAYMENT_BUFFER)
// 	batchChan   = make(chan []models.Payment, BATCH_BUFFER)
// )

// // Publica pagamentos individuais no canal
// func EnqueuePayment(p models.Payment) {
// 	paymentChan <- p
// }

// // Inicializa o pipeline: agregador + workers
// func StartPaymentPipeline() {
// 	log.Printf("Starting pipeline: BATCH_SIZE=%d, WORKERS=%d\n", BATCH_SIZE, WORKER_COUNT)
// 	go startAggregator()
// 	startWorkerPool()
// }

// // Agrega pagamentos em lotes e envia ao canal de batch
// func startAggregator() {
// 	batch := make([]models.Payment, 0, BATCH_SIZE)
// 	ticker := time.NewTicker(TIME_TICKER)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case p := <-paymentChan:
// 			batch = append(batch, p)
// 			if len(batch) >= BATCH_SIZE {
// 				batchChan <- batch
// 				batch = make([]models.Payment, 0, BATCH_SIZE)
// 			}

// 		case <-ticker.C:
// 			if len(batch) > 0 {
// 				batchChan <- batch
// 				batch = make([]models.Payment, 0, BATCH_SIZE)
// 			}
// 		}
// 	}
// }

// // Lança múltiplos workers para inserir batches
// func startWorkerPool() {
// 	for i := range WORKER_COUNT {
// 		go func(id int) {
// 			for batch := range batchChan {
// 				start := time.Now()
// 				if err := models.BulkInsert(batch); err != nil {
// 					log.Printf("[Worker %d] Erro no insert: %v", id, err)
// 				} else {
// 					log.Printf("[Worker %d] Inseriu %d registros em %s", id, len(batch), time.Since(start))
// 				}
// 			}
// 		}(i)
// 	}
// }
