//Custom producer and consumer will have to be finished another time....

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

type check struct {
	partition int32
	offset    int64
	message   string
}

const (
	topic        = "test_topic"
	messageCount = 15
	//clientID      = "test_client"
	numPartitions = int32(8)
)

func main() {
	producer, err := newProducer()
	if err != nil {
		log.Fatal("Unable to create producer", err)
	}

	consumer, err := newConsumer()
	if err != nil {
		log.Fatal("Unable to create consumer", err)
	}

	subscribe(topic, consumer)
	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello from Nolan #%d!", i)
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		})
		if err != nil {
			panic(err)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("Input was: %q\n", line)
		if line == "exit" {
			fmt.Println("Exiting Program...")
			break
		} else {
			_, _, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(line),
			})
			if err != nil {
				panic(err)
			}
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error encountered:", err)
	}
}

//Producer Functions....
func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.ChannelBufferSize = 1
	config.Version = sarama.V0_10_0_1
	config.Producer.Return.Successes = true

	brokers := []string{"127.0.0.1:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)

	return producer, err
}

func sendMessage(topic string, producer sarama.SyncProducer, message string) (int32, int64, error) {
	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})
	return partition, offset, err
}

//Consumer Functions....
func newConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ChannelBufferSize = 1
	config.Version = sarama.V0_10_0_1
	config.Producer.Return.Successes = true

	brokers := []string{"127.0.0.1:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	return consumer, err
}

//Consumer functions
func subscribe(topic string, consumer sarama.Consumer) {
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("Error retrieving partition list")
	}
	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, partition, initialOffset)
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				fmt.Println(message)
			}
		}(pc)
	}
}
