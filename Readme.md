# KV-Store

KV-Store is an in-memory database designed to provide basic key-value store functionality with GET, SET, DEL, and TTL, (time-to-live) capabilities. It's compatible with the Redis protocol, which means any Redis client can connect to the KV-Store.

## Usage

To try out this simple key-value store, follow these steps:

1. **Download the KV-Store binary**:

    ```bash
    $ wget https://github.com/subrag/kv-store/raw/master/kv-store
    ```

2. **Start the KV-Store server on port 8080**:

    ```bash
    $ ./kv-store -p 8080
    ```

    The `-p` flag specifies the port on which the server will listen.

3. **Connect to the KV-Store using the `redis-cli`**. Here, we're connecting to the KV-Store running on host `0.0.0.0` and port `8080`:

    ```bash
    $ redis-cli -h 0.0.0.0 -p 8080
    ```

4. **Interact with the KV-Store** just like you would with Redis. For example:

    ```bash
    0.0.0.0:8080> SET A 10
    OK
    0.0.0.0:8080> GET A
    10
    0.0.0.0:8080> SET B 10 EX 1
    OK
    0.0.0.0:8080> GET B

    0.0.0.0:8080> SET B 10 EX 100
    OK
    0.0.0.0:8080> GET B
    10
    0.0.0.0:8080> 
    ```

## References

- [Understanding epoll](https://medium.com/@avocadi/what-is-epoll-9bbc74272f7c): An excellent resource to understand how the Linux epoll mechanism works, which is crucial for efficient event handling in your KV-Store.

- [Redis Protocol Specification](https://redis.io/docs/reference/protocol-spec/): The KV-Store adheres to the Redis protocol specification, allowing compatibility with Redis clients.

- [Designing a TTL-Based Cache](https://nulpointerexception.com/2023/03/19/design-a-ttl-based-in-memory-cache-in-golang/): This blog post might provide you with insights into TTL functionality in db.

- [Sample Go epoll Implementation](https://gist.github.com/tevino/3a4f4ec4ea9d0ca66d4f): A sample Go epoll implementation that could be beneficial in understanding epoll, creation of fd, how syscall will be using go implement.
