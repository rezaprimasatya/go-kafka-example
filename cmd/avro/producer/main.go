package main

import (
	"fmt"
	"kafka-test/avro/schema"
	"time"

	"github.com/Shopify/sarama"
	"kafka-test/avro/producer"
)

const (
	kafkaHost     = "127.0.0.1:9092"
	kafkaUserName = ""
	kafkaPassword = ""
)

func main() {

	kafkaConfig := initKafkaConfig(kafkaUserName, kafkaPassword)
	producers, err := sarama.NewSyncProducer([]string{kafkaHost}, kafkaConfig)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := producers.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	fmt.Println("Success create kafka sync-producer")

	kafka := &producer.KafkaProducer{
		Producer: producers,
	}

	data := schema.TestSchema{
		Name: "Test",
		Age:  12,
	}

	err = kafka.SendMessage("test-topic-avro", data)
	if err != nil {
		panic(err)
	}
}

func initKafkaConfig(username, password string) *sarama.Config {
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
