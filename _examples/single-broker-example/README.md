# A Go Sarama Example

## 1. Start a single broker node

```bash
$ nolan broker \
          --data-dir="/tmp/nolan0" \
          --broker-addr=127.0.0.1:9092 \
          --raft-addr=127.0.0.1:9093 \
          --serf-addr=127.0.0.1:9094 \
          --bootstrap \
          --bootstrap-expect=1 \
          --id=1
```

## 2. Create a topic for the test node

```bash
$ nolan topic create \
          --broker-addr=127.0.0.1:9092 \
          --partitions=4 \
          --replication-factor=1 \
          --topic=test_topic
```

## 3. Start the sample Sarama script

```bash
$ cd _examples/single-node

$ go run main.go
```
