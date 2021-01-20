
## Start a broker node

```bash
$ nolan broker \
          --data-dir="/tmp/nolan0" \
          --broker-addr=127.0.0.1:9001 \
          --raft-addr=127.0.0.1:9002 \
          --serf-addr=127.0.0.1:9003 \
          --bootstrap \
          --bootstrap-expect=1 \
          --id=1
```

## Create a topic
```bash
$ nolan topic create \
          --broker-addr=127.0.0.1:9001 \
          --partitions=4 \
          --replication-factor=1 \
          --topic=test_topic \
```

