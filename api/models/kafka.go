package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Conn

type Payment struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
	Processor     string  `json:"processor"`
}

func KafkaConnect() {

	var err error

	ctx := context.Background()
	topic := os.Getenv("KAFKA_TOPIC")
	address := os.Getenv("KAFKA_BROKER")
	partition := 0

	KafkaWriter, err = kafka.DialLeader(ctx, "tcp", address, topic, partition)
	if err != nil {
		log.Fatalf("Failed to connect Kafka: %v\n", err)
	}

}

func SendKafka(correlation_id string, amount float64) error {

	payment := Payment{
		CorrelationId: correlation_id,
		Amount:        amount,
	}

	data, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("error marshalling payment: %v", err)
	}

	KafkaWriter.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = KafkaWriter.WriteMessages(
		kafka.Message{Value: data},
	)

	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %v", err)
	}

	return nil

}
