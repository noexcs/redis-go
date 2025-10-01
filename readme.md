> ⚠️ 注意：此项目仅供学习目的，不应用于生产环境。

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
- LIST
  - [x] LPUSH
  - [x] RPUSH
  - [x] LPOP
  - [x] RPOP
- SET
  - [x] SADD
  - [x] SMEMBERS
  - [x] SISMEMBER
  - [x] SREM
- STRING
  - [x] SET
  
    ```sh
    SET key value [NX | XX] [GET] [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | KEEPTTL]
    ```
  - [x] GET
  - [x] INCR
- ZSET
  - [x] ZADD				

# 过期键删除策略

本项目实现了三种过期键删除策略：

1. **定时删除**：在固定时间间隔内，主动遍历并删除所有过期键
2. **惰性删除**：在访问键时检查是否过期，如果过期则删除
3. **定期删除**：定期随机检查部分键，删除其中过期的键

# Benchmark Results

- **OS:** Windows
- **Architecture:** amd64
- **Package:** github.com/noexcs/redis-go/benchmark
- **CPU:** 13th Gen Intel(R) Core(TM) i5-13600KF

## redis-Benchmark

<img src=".\img\p50-latency-multithreaded-vs-singlethreaded.png" style="zoom: 50%" alt="demonstrate"/>

---

<img src=".\img\throughout-multithreaded-vs-singlethreaded.png" style="zoom: 50%" alt="demonstrate"/>

```shell
> neofetch

OS: Ubuntu 24.04.1 LTS on Windows 10 x86_64
Kernel: 5.15.167.4-microsoft-standard-WSL2
Uptime: 1 hour, 42 mins
Packages: 766 (dpkg)
Shell: bash 5.2.21
Theme: Adwaita [GTK3]
Icons: Adwaita [GTK3]
Terminal: Windows Terminal
CPU: 13th Gen Intel i5-13600KF (20) @ 3.494GHz
GPU: 38f5:00:00.0 Microsoft Corporation Basic Render Driver
Memory: 551MiB / 15910MiB

> redis-benchmark -p 6378 -t set,get,incr,hset,hget,lpush,lpop,sadd,smembers,zadd

WARNING: Could not fetch server CONFIG
====== SET ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 90826.52 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.509     0.064     0.431     1.079     1.767     3.863
====== GET ======
  100000 requests completed in 1.07 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 93370.68 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.496     0.064     0.423     1.055     1.647     5.143
====== INCR ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 92421.44 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.501     0.064     0.431     1.031     1.639     4.999
====== LPUSH ======
  100000 requests completed in 1.11 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 90252.70 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.513     0.072     0.439     1.071     1.767     4.663
====== LPOP ======
  100000 requests completed in 1.11 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 90334.23 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.511     0.064     0.431     1.079     1.775     7.975
====== SADD ======
  100000 requests completed in 1.07 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 93720.71 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.489     0.064     0.423     1.031     1.591     3.719
====== HSET ======
  100000 requests completed in 1.09 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 91911.76 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.503     0.072     0.431     1.039     1.655     4.303
====== ZADD ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
...

Cumulative distribution of latencies:
...

Summary:
  throughput summary: 91324.20 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.505     0.072     0.439     1.063     1.639     4.055
```


# 演示

<img src=".\img\demonstrate.png" style="zoom: 50%" alt="demonstrate"/>