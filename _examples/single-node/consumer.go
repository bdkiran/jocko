package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func newConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ChannelBufferSize = 1
	config.Version = sarama.V0_10_0_1
	config.Producer.Return.Successes = true

	brokers := []string{"127.0.0.1:9092"}
	consumer, err := sarama.NewConsumer(brokers, config)
	return consumer, err
}

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
