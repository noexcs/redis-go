start the server:
```shell
go run main.go
```

start the Test:
```shell
ginkgo -r .
```


start the official server Redis 7.0.8 in wsl
```shell
cd /usr/bin
sudo ./redis-server --port 6399
```

# 已完成命令

 - [x] AUTH

   ```shell
   AUTH <password> 
   ```
    前提是在配置文件（默认redis.conf）设置了 `requirepass` 项。
   
 - [x] PING
   
    ```shell
    PING <message>
    ```
 - [x] DEL
    
 - HASH
   - [x] HSET
   - [x] HGET
   - [ ] HDEL
   - [ ] HEXISTS
   - [ ] HGETALL
   - [ ] HINCRBY
   - [ ] HKEYS
   - [ ] HLEN
   - [ ] HVALS
- LIST
  - [x] LPUSH
  - [x] RPUSH
  - [x] LPOP
  - [x] RPOP
- SET
  - [x] SADD
  - [x] SCARD
  - [x] SMEMBERS
  - [x] SISMEMBER
  - [x] SREM
  - [ ] SDIFF
  - [ ] SDIFFSTORE
  - [ ] SINTER
  - [ ] SINTERSTORE
  - [ ] SMOVE
  - [ ] SPOP
  - [ ] SRANDMEMBER
  - [ ] SUNION
  - [ ] SUNIONSTORE
  - [ ] SSCAN
- STRING
  - [x] SET
  - [x] GET
  - [x] GETRANGE
  - [x] INCR
- ZSET
  - [x] ZADD				
  - [ ] ZCARD				
  - [ ] ZCOUNT			
  - [ ] ZINCRBY			
  - [ ] ZINTERSTORE		
  - [ ] ZLEXCOUNT			
  - [ ] ZRANGE			
  - [ ] ZRANGEBYLEX		
  - [ ] ZRANGEBYSCORE		
  - [ ] ZRANK				
  - [ ] ZREM				
  - [ ] ZREMRANGEBYLEX
  - [ ] ZREMRANGEBYRANK
  - [ ] ZREMRANGEBYSCORE
  - [ ] ZREVRANGE			
  - [ ] ZREVRANGEBYSCORE
  - [ ] ZREVRANK			
  - [ ] ZSCORE			
  - [ ] ZUNIONSTORE		
  - [ ] ZSCAN

# Benchmark Results

- **OS:** Windows
- **Architecture:** amd64
- **Package:** github.com/noexcs/redis-go/benchmark
- **CPU:** 13th Gen Intel(R) Core(TM) i5-13600KF

## redis-go Benchmarks

| Benchmark       | Iterations | ns/op     |
|-----------------|------------|-----------|
| BenchmarkRedisSet | 47782      | 24615 ns/op |
| BenchmarkRedisGet | 49344      | 24318 ns/op |
| BenchmarkRedisIncr| 48984      | 24458 ns/op |
| BenchmarkRedisDel | 49016      | 23999 ns/op |


## redis-64-3.0.504 Benchmarks

| Benchmark       | Iterations | ns/op     |
|-----------------|------------|-----------|
| BenchmarkRedisSet | 50744      | 22828 ns/op |
| BenchmarkRedisGet | 52902      | 22471 ns/op |
| BenchmarkRedisIncr| 52924      | 22658 ns/op |
| BenchmarkRedisDel | 52605      | 22439 ns/op |

# 演示

<img src=".\img\demonstrate.png" style="zoom: 50%" alt="demonstrate"/>
