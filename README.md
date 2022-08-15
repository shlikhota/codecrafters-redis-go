This is a Go solutions to the
["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

A toy Redis clone that's capable of handling
basic commands like `PING`, `SET` and `GET`.

# How to run it

It starts a local server on port 6379:

```sh
./spawn_redis_server.sh
```

Let's check the functionality:
```sh
~ redis-cli
127.0.0.1:6379> GET somekey
(nil)
127.0.0.1:6379> SET somekey somevalue
OK
127.0.0.1:6379> GET somekey
"somevalue"
127.0.0.1:6379> SET somekey somevalue PX 3000
OK
127.0.0.1:6379> GET somekey
"somevalue"
127.0.0.1:6379> GET somekey
(nil)
```
