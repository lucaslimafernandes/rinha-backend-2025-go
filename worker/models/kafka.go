package models

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaReader *kafka.Reader

// type Payment struct {
// 	CorrelationId string  `json:"correlationId"`
// 	Amount        float64 `json:"amount"`
// 	Processor     string  `json:"processor"`
// }

func KafkaConnect() {
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // tempo para commit automático
	})
	log.Println("Conectado ao Kafka como consumidor.")
}

func ReadKafka(ctx context.Context, handler func([]byte) error) {
	defer KafkaReader.Close()

	log.Println("Iniciando leitura de mensagens Kafka...")

	for {
		m, err := KafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			continue
		}

		log.Printf("Mensagem recebida [%s]: %s", string(m.Key), string(m.Value))

		if err := handler(m.Value); err != nil {
			log.Printf("Erro ao processar mensagem: %v", err)
			// você pode decidir salvar em uma DLQ aqui
		}
	}
}
