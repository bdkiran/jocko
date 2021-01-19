package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cli = &cobra.Command{
		Use:   "nolan",
		Short: "Kafka in Go",
		Long: `----------------------
        nolan 
---------------------
A distributed commit log built in go.
No messy JVM or complex configurations.`,
	}

	// brokerCfg = config.DefaultConfig()

	// topicCfg = struct {
	// 	BrokerAddr        string
	// 	Topic             string
	// 	Partitions        int32
	// 	ReplicationFactor int
	// }{}
)

//Execute executes the cli command
func Execute() error {
	return cli.Execute()
}

func init() {
	cli.AddCommand(brokerCmd)
	cli.AddCommand(topicCmd)
}
