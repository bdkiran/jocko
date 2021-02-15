# nolan cluster example

This will start a local three node cluster.

## Build

```bash
$ cd /nolan

$ go build ./...

$ go install
```

## Start multinode clusters

```bash
$ nolan broker \
          --data-dir="/tmp/nolan0" \
          --broker-addr=127.0.0.1:9001 \
          --raft-addr=127.0.0.1:9002 \
          --serf-addr=127.0.0.1:9003 \
          --bootstrap \
          --bootstrap-expect=3 \
          --id=1

$ nolan broker \
          --data-dir="/tmp/nolan1" \
          --broker-addr=127.0.0.1:9101 \
          --raft-addr=127.0.0.1:9102 \
          --serf-addr=127.0.0.1:9103 \
          --join=127.0.0.1:9003 \
          --join-wan=127.0.0.1:9003 \
          --bootstrap-expect=3 \
          --id=2

$ nolan broker \
          --data-dir="/tmp/nolan2" \
          --broker-addr=127.0.0.1:920 \
          --raft-addr=127.0.0.1:9202 \
          --serf-addr=127.0.0.1:9203 \
          --join=127.0.0.1:9003 \
          --bootstrap-expect=3 \
          --id=3
```

## docker-compose cluster

To start a [docker compose](https://docs.docker.com/compose/) cluster use the provided `docker-compose.yml`.
