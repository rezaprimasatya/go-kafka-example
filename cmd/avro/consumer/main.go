package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"kafka-test/avro/consumer"

	"github.com/Shopify/sarama"
)

const (
	kafkaHost     = "127.0.0.1:9092"
	kafkaUserName = ""
	kafkaPassword = ""
)

func main() {

	kafkaConfig := getKafkaConfig(kafkaUserName, kafkaPassword)

	consumers, err := sarama.NewConsumer([]string{kafkaHost}, kafkaConfig)
	if err != nil {
		fmt.Println( err)
	}
	defer func() {
		if err := consumers.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}()

	kafkaConsumer := &consumer.KafkaConsumer{
		Consumer: consumers,
	}

	signals := make(chan os.Signal, 1)
	kafkaConsumer.Consume([]string{"test-topic-avro"}, signals)
}

func getKafkaConfig(username, password string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}