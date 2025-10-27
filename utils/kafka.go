package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func SendKafkaMessage(topic string, data interface{}) {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Println("Kafka connection failed:", err)
		return
	}
	defer producer.Close()

	jsonData, _ := json.Marshal(data)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonData),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Println("Kafka send error:", err)
		return
	}

	fmt.Println("📤 Kafka log terkirim ke topic:", topic)
}
