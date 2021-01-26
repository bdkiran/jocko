# Nolan

![ci](https://github.com/bdkiran/nolan/workflows/Go/badge.svg)

*nolan* is a lightweight distributed a commit log service. Create clusters trivally: straightforward and not a ton of configuation needed. Avalible to run as a standalone binary or in containers.

This is a fork of the Jocko Repository.

## Why?
*nolan* is intended to reduce overhead that comes with setting up a distributed commit log or queue. Weather it's used for website activity tracking or application log aggregation, this tool should be able to provide stable groundwork for many use cases. 

## How?
*nolan* follows design patterns very similar to Kafka. Without a zookeeper dependency, nolan utilizes [raft](https://raft.github.io/) to achieve its distributed and scalable goals. 

## Gaols of the Project

- Implement Kafka in Go
- Protocol compatible with Kafka so Kafka clients and services work with Noloan
- Make operating simpler
- Distribute a single binary
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

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## License

Nolan is under the MIT license, see the [LICENSE](LICENSE) file for details.

### Further Reading

- [How Jocko's built-in service discovery and consensus works](https://medium.com/the-hoard/building-a-kafka-that-doesnt-depend-on-zookeeper-2c4701b6e961#.uamxtq1yz)
- [How Jocko's (and Kafka's) storage internals work](https://medium.com/the-hoard/how-kafkas-storage-internals-work-3a29b02e026#.qfbssm978)
- [J Kreps: The Log](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying)
