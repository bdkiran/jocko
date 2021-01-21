
## Start a broker node

$ nolan broker test

## Create a topic
```bash
$ nolan topic create \
          --broker-addr=127.0.0.1:9092 \
          --partitions=4 \
          --replication-factor=1 \
          --topic=test_topic
```

