package producer

import (
	"bytes"
	"fmt"
	"github.com/Shopify/sarama"
	"kafka-test/avro/schema"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}


func (p *KafkaProducer) SendMessage(topic string, data schema.TestSchema) error {

	var buf bytes.Buffer
	fmt.Printf("Serializing struct: %#v\n", data)
	data.Serialize(&buf)

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(buf.Bytes()),
	}

	partition, offset, err := p.Producer.SendMessage(kafkaMsg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Send message success, Topic %v, Partition %v, Offset %d\n", topic, partition, offset)
	return nil
}
