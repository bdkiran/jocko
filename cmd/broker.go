package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	gracefully "github.com/tj/go-gracefully"
	"github.com/travisjeffery/jocko/jocko"
	"github.com/travisjeffery/jocko/jocko/config"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/hashicorp/memberlist"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var (
	brokerCmd = &cobra.Command{
		Use:   "broker",
		Short: "Manage broker",
		Long:  `Everything nolan related with broker. Manage here`,
		Run:   run,
	}

	brokerCfg = config.DefaultConfig()
)

func init() {
	//brokerCmd := &cobra.Command{Use: "broker", Short: "Run a Jocko broker", Run: run, Args: cobra.NoArgs}
	brokerCmd.Flags().StringVar(&brokerCfg.RaftAddr, "raft-addr", "127.0.0.1:9093", "Address for Raft to bind and advertise on")
	brokerCmd.Flags().StringVar(&brokerCfg.DataDir, "data-dir", "/tmp/jocko", "A comma separated list of directories under which to store log files")
	brokerCmd.Flags().StringVar(&brokerCfg.Addr, "broker-addr", "0.0.0.0:9092", "Address for broker to bind on")
	brokerCmd.Flags().Var(newMemberlistConfigValue(brokerCfg.SerfLANConfig.MemberlistConfig, "0.0.0.0:9094"), "serf-addr", "Address for Serf to bind on")
	brokerCmd.Flags().BoolVar(&brokerCfg.Bootstrap, "bootstrap", false, "Initial cluster bootstrap (dangerous!)")
	brokerCmd.Flags().IntVar(&brokerCfg.BootstrapExpect, "bootstrap-expect", 0, "Expected number of nodes in cluster")
	brokerCmd.Flags().StringSliceVar(&brokerCfg.StartJoinAddrsLAN, "join", nil, "Address of an broker serf to join at start time. Can be specified multiple times.")
	brokerCmd.Flags().StringSliceVar(&brokerCfg.StartJoinAddrsWAN, "join-wan", nil, "Address of an broker serf to join -wan at start time. Can be specified multiple times.")
	brokerCmd.Flags().Int32Var(&brokerCfg.ID, "id", 0, "Broker ID")

	//cli.AddCommand(brokerCmd)
}

func run(cmd *cobra.Command, args []string) {
	var err error

	log.SetPrefix(fmt.Sprintf("jocko: node id: %d: ", brokerCfg.ID))

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.New(
		"jocko",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(err)
	}

	broker, err := jocko.NewBroker(brokerCfg, tracer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting broker: %v\n", err)
		os.Exit(1)
	}

	srv := jocko.NewServer(brokerCfg, broker, nil, tracer, closer.Close)
	if err := srv.Start(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}

	defer srv.Shutdown()

	gracefully.Timeout = 10 * time.Second
	gracefully.Shutdown()

	if err := broker.Shutdown(); err != nil {
		fmt.Fprintf(os.Stderr, "error shutting down store: %v\n", err)
		os.Exit(1)
	}
}

type memberlistConfigValue memberlist.Config

func newMemberlistConfigValue(p *memberlist.Config, val string) (m *memberlistConfigValue) {
	m = (*memberlistConfigValue)(p)
	m.Set(val)
	return
}

func (v *memberlistConfigValue) Set(s string) error {
	bindIP, bindPort, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	v.BindAddr = bindIP
	v.BindPort, err = strconv.Atoi(bindPort)
	if err != nil {
		return err
	}
	return nil
}

func (v *memberlistConfigValue) Type() string {
	return "string"
}

func (v *memberlistConfigValue) String() string {
	return fmt.Sprintf("%s:%d", v.BindAddr, v.BindPort)
}
