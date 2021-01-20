package main

import "github.com/Shopify/sarama"

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
