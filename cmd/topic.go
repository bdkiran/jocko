package cmd

import (
	"fmt"
	"os"

	"github.com/bdkiran/nolan/nolan"
	"github.com/bdkiran/nolan/protocol"
	"github.com/spf13/cobra"
)

var (
	topicCmd = &cobra.Command{
		Use:   "topic",
		Short: "Manage topics",
		Long:  `Everything about nolan topics here. Manage them here`,
	}

	topicCfg = struct {
		BrokerAddr        string
		Topic             string
		Partitions        int32
		ReplicationFactor int
	}{}
)

func init() {
	//Create a new command to 'create topic'
	createTopicCmd := &cobra.Command{Use: "create", Short: "Create a topic", Run: createTopic, Args: cobra.NoArgs}
	createTopicCmd.Flags().StringVar(&topicCfg.BrokerAddr, "broker-addr", "0.0.0.0:9092", "Address for Broker to bind on")
	createTopicCmd.Flags().StringVar(&topicCfg.Topic, "topic", "", "Name of topic to create (required)")
	createTopicCmd.MarkFlagRequired("topic")
	createTopicCmd.Flags().Int32Var(&topicCfg.Partitions, "partitions", 1, "Number of partitions")
	createTopicCmd.Flags().IntVar(&topicCfg.ReplicationFactor, "replication-factor", 1, "Replication factor")
	topicCmd.AddCommand(createTopicCmd)
}

func createTopic(cmd *cobra.Command, args []string) {
	conn, err := nolan.Dial("tcp", topicCfg.BrokerAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to broker: %v\n", err)
		os.Exit(1)
	}

	resp, err := conn.CreateTopics(&protocol.CreateTopicRequests{
		Requests: []*protocol.CreateTopicRequest{{
			Topic:             topicCfg.Topic,
			NumPartitions:     topicCfg.Partitions,
			ReplicationFactor: int16(topicCfg.ReplicationFactor),
			ReplicaAssignment: nil,
			Configs:           nil,
		}},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with request to broker: %v\n", err)
		os.Exit(1)
	}
	for _, topicErrCode := range resp.TopicErrorCodes {
		if topicErrCode.ErrorCode != protocol.ErrNone.Code() {
			err := protocol.Errs[topicErrCode.ErrorCode]
			fmt.Fprintf(os.Stderr, "error code: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("created topic: %v\n", topicCfg.Topic)
}
