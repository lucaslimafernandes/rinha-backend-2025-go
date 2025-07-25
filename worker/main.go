package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
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

func handler(payment models.Payment) {

	statusMutex.RLock()
	defer statusMutex.RUnlock()

	useDefault := healthStatus["default"] >= 0 && healthStatus["default"] < 99
	useFallback := healthStatus["fallback"] >= 0

	switch {
	case useDefault:
		err := service.PaymentSend("default", payment)
		if err != nil {
			log.Println("Internal Server Error")
		}
	case !useDefault && useFallback:
		err := service.PaymentSend("fallback", payment)
		if err != nil {
			log.Println("Internal Server Error")
		}

	default:
		err := service.PaymentSend("default", payment)
		if err != nil {
			log.Println("Internal Server Error")
		}
	}

}

func main() {

	// PgBouncer
	os.Setenv("PGAPPNAME", "")
	os.Setenv("PGOPTIONS", "")

	models.KafkaConnect()

	go updateHealthLoop()

	hand := func(msg []byte) error {
		log.Printf("processing: %s\n", string(msg))

		var payment models.Payment

		err := json.Unmarshal(msg, &payment)
		if err != nil {
			log.Printf("error to unmarshall: %v\n", err)
			return err
		}

		handler(payment)
		return nil
	}

	ctx := context.Background()
	models.ReadKafka(ctx, hand)

}
