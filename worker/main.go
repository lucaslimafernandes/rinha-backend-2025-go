package main

import (
	"context"
	"log"
	"sync"
	"time"
	"worker-rinha/models"
	"worker-rinha/service"
)

type HealthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var (
	healthStatus = make(map[string]int)
	statusMutex  = sync.RWMutex{}
)

func updateHealthLoop() {

	for {
		ch := make(chan map[string]int)

		go service.CheckService(ch)
		result := <-ch

		statusMutex.Lock()
		healthStatus = result
		statusMutex.Unlock()

		time.Sleep(5 * time.Second)

	}
}

// func aux() {

// 	statusMutex.RLock()
// 	defer statusMutex.RUnlock()

// 	if healthStatus["default"] >= 0 {
// 		fmt.Fprintf(w, "default")
// 	}

// 	switch {
// 	case healthStatus["default"] >= 0 && healthStatus["default"] < 99:
// 		err := utils.PaymentSend("default", payment)
// 		if err != nil {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}
// 	case healthStatus["default"] == -1 && healthStatus["fallback"] >= 0:
// 		err = utils.PaymentSend("fallback", payment)
// 		if err != nil {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}

// 	case healthStatus["default"] > 99 && healthStatus["fallback"] >= 0:
// 		err = utils.PaymentSend("fallback", payment)
// 		if err != nil {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}

// 	default:
// 		err = utils.PaymentSend("default", payment)
// 		if err != nil {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}
// 	}

// }

func main() {

	models.KafkaConnect()

	go updateHealthLoop()

	hand := func(msg []byte) error {
		log.Printf("processing: %s\n", string(msg))

		return nil
	}

	ctx := context.Background()
	models.ReadKafka(ctx, hand)

}
