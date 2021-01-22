# Nolan

![ci](https://github.com/bdkiran/nolan/workflows/Go/badge.svg)

A distributed commit log/write ahead log service written in in Go.

This is a fork of the Jocko Repository.

## Gaols of the Project

- Implement Kafka in Go
- Protocol compatible with Kafka so Kafka clients and services work with Jocko
- Make operating simpler
- Distribute a single binary
- Use Serf for discovery, Raft for consensus (and remove the need to run ZooKeeper)
- Simpler configuration settings
  - Get a cluster or single broker up and running quickly

## Building

### Local

1. Clone Nolan

    ```bash
    $ git clone github.com/bdkiran/nolan
    ```

2. Build Nolan

    ```bash
    $ cd nolan
    $ go build ./...
    $ go install
    ```

3. Try running one of the examples

### Docker

`docker build -t travisjeffery/jocko:latest .`

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## License

Nolan is under the MIT license, see the [LICENSE](LICENSE) file for details.

### Further Reading

- [How Jocko's built-in service discovery and consensus works](https://medium.com/the-hoard/building-a-kafka-that-doesnt-depend-on-zookeeper-2c4701b6e961#.uamxtq1yz)
- [How Jocko's (and Kafka's) storage internals work](https://medium.com/the-hoard/how-kafkas-storage-internals-work-3a29b02e026#.qfbssm978)
- [J Kreps: The Log](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying)
