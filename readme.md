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