package producer

import (
	"fmt"
	"github.com/Shopify/sarama"
)


type KafkaProducer struct {
	Producer sarama.SyncProducer
}


func (p *KafkaProducer) SendMessage(topic, msg string) error {

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := p.Producer.SendMessage(kafkaMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Send message success, Topic %v, Partition %v, Offset %d\n", topic, partition, offset)
	return nil
}
