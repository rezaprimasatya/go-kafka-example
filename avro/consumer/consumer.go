package consumer

import (
	"bytes"
	"fmt"
	"github.com/Shopify/sarama"
	"kafka-test/avro/schema"
	"os"
)

type KafkaConsumer struct {
	Consumer sarama.Consumer
}

func (c *KafkaConsumer) Consume(topics []string, signals chan os.Signal) {
	chanMessage := make(chan *sarama.ConsumerMessage, 256)

	for _, topic := range topics {
		partitionList, err := c.Consumer.Partitions(topic)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(c.Consumer, topic, partition, chanMessage)
		}
	}
	fmt.Println("Kafka is consuming....")

	// https://www.ardanlabs.com/blog/2013/11/label-breaks-in-go.html
ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			buf := bytes.NewBuffer(msg.Value)
			testStruct, err := schema.DeserializeTestSchema(buf)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("New Message from kafka, message [name]: %v\n", testStruct.Name)
			fmt.Printf("New Message from kafka, message [age]: %v\n", testStruct.Age)
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		fmt.Printf("Unable to consume partition %v got error %v\n", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			fmt.Printf("Unable to close partition %v: %v\n", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}

}
