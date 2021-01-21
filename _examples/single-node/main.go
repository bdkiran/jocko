package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/bdkiran/nolan/log"
)

type check struct {
	partition int32
	offset    int64
	message   string
}

const (
	topic         = "test_topic"
	messageCount  = 15
	clientID      = "test_client"
	numPartitions = int32(8)
)

func main() {
	config := sarama.NewConfig()
	config.ChannelBufferSize = 1
	config.Version = sarama.V0_10_0_1
	config.Producer.Return.Successes = true

	brokers := []string{"127.0.0.1:9092"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	pmap := make(map[int32][]check)

	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello from Jocko #%d!", i)
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		})
		if err != nil {
			panic(err)
		}
		pmap[partition] = append(pmap[partition], check{
			partition: partition,
			offset:    offset,
			message:   message,
		})
	}
	if err = producer.Close(); err != nil {
		panic(err)
	}

	var totalChecked int
	for partitionID := range pmap {
		checked := 0
		consumer, err := sarama.NewConsumer(brokers, config)
		if err != nil {
			panic(err)
		}
		partition, err := consumer.ConsumePartition(topic, partitionID, 0)
		if err != nil {
			panic(err)
		}
		i := 0
		for msg := range partition.Messages() {
			log.Info.Println("-----------Consuming new message---------------------")
			log.Info.Printf("msg partition [%d] offset [%d]\n", msg.Partition, msg.Offset)
			check := pmap[partitionID][i]
			if string(msg.Value) != check.message {
				log.Error.Fatalf("msg values not equal: partition: %d: offset: %d", msg.Partition, msg.Offset)
			}
			if msg.Offset != check.offset {
				log.Error.Fatalf("msg offsets not equal: partition: %d: offset: %d", msg.Partition, msg.Offset)
			}
			log.Info.Printf("msg is ok: partition: %d: offset: %d", msg.Partition, msg.Offset)
			log.Info.Println(string(msg.Value))
			i++
			checked++
			fmt.Printf("i: %d, len: %d\n", i, len(pmap[partitionID]))
			if i == len(pmap[partitionID]) {
				totalChecked += checked
				fmt.Println("checked partition:", partitionID)
				if err = consumer.Close(); err != nil {
					panic(err)
				}
				break
			} else {
				fmt.Println("still checking partition:", partitionID)
			}
		}
	}
	log.Info.Println("===========================================================")
	log.Info.Printf("producer and consumer worked! %d messages ok\n", totalChecked)
}

//Custom producer and consumer will have to be finished another time....

// package main

// import (
// 	"fmt"

// 	"github.com/Shopify/sarama"
// )

// type check struct {
// 	partition int32
// 	offset    int64
// 	message   string
// }

// const (
// 	topic        = "test_topic"
// 	messageCount = 15
// 	//clientID      = "test_client"
// 	numPartitions = int32(8)
// )

// func main() {
// 	producer, err := newProducer()
// 	if err != nil {
// 		fmt.Println("Unable to create producer", err)
// 	}

// 	consumer, err := newConsumer()
// 	if err != nil {
// 		fmt.Println("Unable to create consumer", err)
// 	}

// 	subscribe(topic, consumer)
// 	for i := 0; i < messageCount; i++ {
// 		message := fmt.Sprintf("Hello from Nolan #%d!", i)
// 		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
// 			Topic: topic,
// 			Value: sarama.StringEncoder(message),
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// scanner := bufio.NewScanner(os.Stdin)
// for scanner.Scan() {
// 	line := scanner.Text()
// 	//fmt.Printf("Input was: %q\n", line)
// 	if line == "exit" {
// 		fmt.Println("Exiting Program...")
// 		break
// 	} else {
// 		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
// 			Topic: topic,
// 			Value: sarama.StringEncoder(line),
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// }
// if err := scanner.Err(); err != nil {
// 	fmt.Println("Error encountered:", err)
// }
//}

//Producer Functions....
// func newProducer() (sarama.SyncProducer, error) {
// 	config := sarama.NewConfig()
// 	config.ChannelBufferSize = 1
// 	config.Version = sarama.V0_10_0_1
// 	config.Producer.Return.Successes = true

// 	brokers := []string{"127.0.0.1:9092"}
// 	producer, err := sarama.NewSyncProducer(brokers, config)

// 	return producer, err
// }

// func sendMessage(topic string, producer sarama.SyncProducer, message string) (int32, int64, error) {
// 	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(message),
// 	})
// 	return partition, offset, err
// }

// //Consumer Functions....
// func newConsumer() (sarama.Consumer, error) {
// 	config := sarama.NewConfig()
// 	config.ChannelBufferSize = 1
// 	config.Version = sarama.V0_10_0_1
// 	config.Producer.Return.Successes = true

// 	brokers := []string{"127.0.0.1:9092"}
// 	consumer, err := sarama.NewConsumer(brokers, config)
// 	return consumer, err
// }

// func subscribe(topic string, consumer sarama.Consumer) {
// 	partitionList, err := consumer.Partitions(topic)
// 	if err != nil {
// 		fmt.Println("Error retrieving partition list")
// 	}
// 	initialOffset := sarama.OffsetOldest

// 	for _, partition := range partitionList {
// 		pc, _ := consumer.ConsumePartition(topic, partition, initialOffset)
// 		go func(pc sarama.PartitionConsumer) {
// 			for message := range pc.Messages() {
// 				fmt.Println(message)
// 			}
// 		}(pc)
// 	}
// }
