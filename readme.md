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
  100000 requests completed in 1.32 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.103 milliseconds (cumulative count 1)
50.000% <= 0.543 milliseconds (cumulative count 50666)
75.000% <= 0.767 milliseconds (cumulative count 75017)
87.500% <= 0.991 milliseconds (cumulative count 87792)
93.750% <= 1.215 milliseconds (cumulative count 93776)
96.875% <= 1.455 milliseconds (cumulative count 96877)
98.438% <= 1.695 milliseconds (cumulative count 98454)
99.219% <= 1.935 milliseconds (cumulative count 99239)
99.609% <= 2.191 milliseconds (cumulative count 99610)
99.805% <= 2.535 milliseconds (cumulative count 99810)
99.902% <= 2.879 milliseconds (cumulative count 99905)
99.951% <= 3.303 milliseconds (cumulative count 99952)
99.976% <= 3.855 milliseconds (cumulative count 99976)
99.988% <= 4.159 milliseconds (cumulative count 99988)
99.994% <= 4.271 milliseconds (cumulative count 99994)
99.997% <= 4.415 milliseconds (cumulative count 99997)
99.998% <= 4.783 milliseconds (cumulative count 99999)
99.999% <= 5.311 milliseconds (cumulative count 100000)
100.000% <= 5.311 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.001% <= 0.103 milliseconds (cumulative count 1)
2.029% <= 0.207 milliseconds (cumulative count 2029)
12.739% <= 0.303 milliseconds (cumulative count 12739)
29.663% <= 0.407 milliseconds (cumulative count 29663)
44.904% <= 0.503 milliseconds (cumulative count 44904)
59.042% <= 0.607 milliseconds (cumulative count 59042)
69.610% <= 0.703 milliseconds (cumulative count 69610)
78.051% <= 0.807 milliseconds (cumulative count 78051)
83.984% <= 0.903 milliseconds (cumulative count 83984)
88.343% <= 1.007 milliseconds (cumulative count 88343)
91.309% <= 1.103 milliseconds (cumulative count 91309)
93.653% <= 1.207 milliseconds (cumulative count 93653)
95.141% <= 1.303 milliseconds (cumulative count 95141)
96.435% <= 1.407 milliseconds (cumulative count 96435)
97.277% <= 1.503 milliseconds (cumulative count 97277)
97.993% <= 1.607 milliseconds (cumulative count 97993)
98.502% <= 1.703 milliseconds (cumulative count 98502)
98.899% <= 1.807 milliseconds (cumulative count 98899)
99.167% <= 1.903 milliseconds (cumulative count 99167)
99.389% <= 2.007 milliseconds (cumulative count 99389)
99.507% <= 2.103 milliseconds (cumulative count 99507)
99.945% <= 3.103 milliseconds (cumulative count 99945)
99.983% <= 4.103 milliseconds (cumulative count 99983)
99.999% <= 5.103 milliseconds (cumulative count 99999)
100.000% <= 6.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 75930.14 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.623     0.096     0.543     1.295     1.839     5.311
====== GET ======
  100000 requests completed in 1.30 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.087 milliseconds (cumulative count 1)
50.000% <= 0.535 milliseconds (cumulative count 50610)
75.000% <= 0.759 milliseconds (cumulative count 75363)
87.500% <= 0.975 milliseconds (cumulative count 87785)
93.750% <= 1.183 milliseconds (cumulative count 93751)
96.875% <= 1.415 milliseconds (cumulative count 96935)
98.438% <= 1.647 milliseconds (cumulative count 98462)
99.219% <= 1.951 milliseconds (cumulative count 99224)
99.609% <= 2.255 milliseconds (cumulative count 99612)
99.805% <= 2.703 milliseconds (cumulative count 99805)
99.902% <= 2.951 milliseconds (cumulative count 99904)
99.951% <= 3.159 milliseconds (cumulative count 99953)
99.976% <= 3.303 milliseconds (cumulative count 99976)
99.988% <= 3.823 milliseconds (cumulative count 99990)
99.994% <= 4.015 milliseconds (cumulative count 99996)
99.997% <= 4.135 milliseconds (cumulative count 99997)
99.998% <= 4.191 milliseconds (cumulative count 100000)
100.000% <= 4.191 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.004% <= 0.103 milliseconds (cumulative count 4)
2.244% <= 0.207 milliseconds (cumulative count 2244)
13.301% <= 0.303 milliseconds (cumulative count 13301)
30.558% <= 0.407 milliseconds (cumulative count 30558)
46.018% <= 0.503 milliseconds (cumulative count 46018)
60.091% <= 0.607 milliseconds (cumulative count 60091)
70.510% <= 0.703 milliseconds (cumulative count 70510)
78.867% <= 0.807 milliseconds (cumulative count 78867)
84.580% <= 0.903 milliseconds (cumulative count 84580)
88.917% <= 1.007 milliseconds (cumulative count 88917)
91.870% <= 1.103 milliseconds (cumulative count 91870)
94.181% <= 1.207 milliseconds (cumulative count 94181)
95.665% <= 1.303 milliseconds (cumulative count 95665)
96.857% <= 1.407 milliseconds (cumulative count 96857)
97.643% <= 1.503 milliseconds (cumulative count 97643)
98.285% <= 1.607 milliseconds (cumulative count 98285)
98.676% <= 1.703 milliseconds (cumulative count 98676)
98.961% <= 1.807 milliseconds (cumulative count 98961)
99.142% <= 1.903 milliseconds (cumulative count 99142)
99.314% <= 2.007 milliseconds (cumulative count 99314)
99.450% <= 2.103 milliseconds (cumulative count 99450)
99.938% <= 3.103 milliseconds (cumulative count 99938)
99.996% <= 4.103 milliseconds (cumulative count 99996)
100.000% <= 5.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 77041.60 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.613     0.080     0.535     1.263     1.831     4.191
====== INCR ======
  100000 requests completed in 1.30 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.095 milliseconds (cumulative count 2)
50.000% <= 0.535 milliseconds (cumulative count 50542)
75.000% <= 0.759 milliseconds (cumulative count 75325)
87.500% <= 0.967 milliseconds (cumulative count 87740)
93.750% <= 1.183 milliseconds (cumulative count 93783)
96.875% <= 1.423 milliseconds (cumulative count 96904)
98.438% <= 1.679 milliseconds (cumulative count 98462)
99.219% <= 1.895 milliseconds (cumulative count 99225)
99.609% <= 2.111 milliseconds (cumulative count 99618)
99.805% <= 2.335 milliseconds (cumulative count 99805)
99.902% <= 2.575 milliseconds (cumulative count 99904)
99.951% <= 2.871 milliseconds (cumulative count 99952)
99.976% <= 3.095 milliseconds (cumulative count 99976)
99.988% <= 3.247 milliseconds (cumulative count 99988)
99.994% <= 3.543 milliseconds (cumulative count 99994)
99.997% <= 3.639 milliseconds (cumulative count 99997)
99.998% <= 3.727 milliseconds (cumulative count 99999)
99.999% <= 3.751 milliseconds (cumulative count 100000)
100.000% <= 3.751 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.003% <= 0.103 milliseconds (cumulative count 3)
2.118% <= 0.207 milliseconds (cumulative count 2118)
13.041% <= 0.303 milliseconds (cumulative count 13041)
30.208% <= 0.407 milliseconds (cumulative count 30208)
45.790% <= 0.503 milliseconds (cumulative count 45790)
60.066% <= 0.607 milliseconds (cumulative count 60066)
70.466% <= 0.703 milliseconds (cumulative count 70466)
78.980% <= 0.807 milliseconds (cumulative count 78980)
84.907% <= 0.903 milliseconds (cumulative count 84907)
89.254% <= 1.007 milliseconds (cumulative count 89254)
92.106% <= 1.103 milliseconds (cumulative count 92106)
94.255% <= 1.207 milliseconds (cumulative count 94255)
95.630% <= 1.303 milliseconds (cumulative count 95630)
96.776% <= 1.407 milliseconds (cumulative count 96776)
97.485% <= 1.503 milliseconds (cumulative count 97485)
98.086% <= 1.607 milliseconds (cumulative count 98086)
98.571% <= 1.703 milliseconds (cumulative count 98571)
98.955% <= 1.807 milliseconds (cumulative count 98955)
99.239% <= 1.903 milliseconds (cumulative count 99239)
99.458% <= 2.007 milliseconds (cumulative count 99458)
99.602% <= 2.103 milliseconds (cumulative count 99602)
99.976% <= 3.103 milliseconds (cumulative count 99976)
100.000% <= 4.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 77041.60 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.612     0.088     0.535     1.255     1.823     3.751
====== LPUSH ======
  100000 requests completed in 1.35 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.087 milliseconds (cumulative count 2)
50.000% <= 0.543 milliseconds (cumulative count 50109)
75.000% <= 0.783 milliseconds (cumulative count 75529)
87.500% <= 0.999 milliseconds (cumulative count 87573)
93.750% <= 1.231 milliseconds (cumulative count 93802)
96.875% <= 1.519 milliseconds (cumulative count 96884)
98.438% <= 1.863 milliseconds (cumulative count 98466)
99.219% <= 2.215 milliseconds (cumulative count 99224)
99.609% <= 2.575 milliseconds (cumulative count 99615)
99.805% <= 2.911 milliseconds (cumulative count 99806)
99.902% <= 3.311 milliseconds (cumulative count 99905)
99.951% <= 3.703 milliseconds (cumulative count 99953)
99.976% <= 3.927 milliseconds (cumulative count 99976)
99.988% <= 4.063 milliseconds (cumulative count 99988)
99.994% <= 4.199 milliseconds (cumulative count 99994)
99.997% <= 4.607 milliseconds (cumulative count 99997)
99.998% <= 4.775 milliseconds (cumulative count 99999)
99.999% <= 4.991 milliseconds (cumulative count 100000)
100.000% <= 4.991 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.006% <= 0.103 milliseconds (cumulative count 6)
2.069% <= 0.207 milliseconds (cumulative count 2069)
12.454% <= 0.303 milliseconds (cumulative count 12454)
29.002% <= 0.407 milliseconds (cumulative count 29002)
44.285% <= 0.503 milliseconds (cumulative count 44285)
58.422% <= 0.607 milliseconds (cumulative count 58422)
68.638% <= 0.703 milliseconds (cumulative count 68638)
77.304% <= 0.807 milliseconds (cumulative count 77304)
83.166% <= 0.903 milliseconds (cumulative count 83166)
87.864% <= 1.007 milliseconds (cumulative count 87864)
90.997% <= 1.103 milliseconds (cumulative count 90997)
93.382% <= 1.207 milliseconds (cumulative count 93382)
94.882% <= 1.303 milliseconds (cumulative count 94882)
96.043% <= 1.407 milliseconds (cumulative count 96043)
96.773% <= 1.503 milliseconds (cumulative count 96773)
97.427% <= 1.607 milliseconds (cumulative count 97427)
97.889% <= 1.703 milliseconds (cumulative count 97889)
98.291% <= 1.807 milliseconds (cumulative count 98291)
98.567% <= 1.903 milliseconds (cumulative count 98567)
98.813% <= 2.007 milliseconds (cumulative count 98813)
99.007% <= 2.103 milliseconds (cumulative count 99007)
99.862% <= 3.103 milliseconds (cumulative count 99862)
99.989% <= 4.103 milliseconds (cumulative count 99989)
100.000% <= 5.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 74294.21 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.635     0.080     0.543     1.319     2.103     4.991
====== LPOP ======
  100000 requests completed in 1.33 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.095 milliseconds (cumulative count 2)
50.000% <= 0.543 milliseconds (cumulative count 50569)
75.000% <= 0.775 milliseconds (cumulative count 75205)
87.500% <= 0.999 milliseconds (cumulative count 87585)
93.750% <= 1.239 milliseconds (cumulative count 93876)
96.875% <= 1.519 milliseconds (cumulative count 96906)
98.438% <= 1.823 milliseconds (cumulative count 98449)
99.219% <= 2.151 milliseconds (cumulative count 99220)
99.609% <= 2.479 milliseconds (cumulative count 99612)
99.805% <= 2.815 milliseconds (cumulative count 99805)
99.902% <= 3.191 milliseconds (cumulative count 99903)
99.951% <= 3.407 milliseconds (cumulative count 99952)
99.976% <= 3.567 milliseconds (cumulative count 99976)
99.988% <= 3.863 milliseconds (cumulative count 99988)
99.994% <= 4.111 milliseconds (cumulative count 99994)
99.997% <= 4.255 milliseconds (cumulative count 99997)
99.998% <= 5.183 milliseconds (cumulative count 99999)
99.999% <= 5.207 milliseconds (cumulative count 100000)
100.000% <= 5.207 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.003% <= 0.103 milliseconds (cumulative count 3)
2.152% <= 0.207 milliseconds (cumulative count 2152)
13.118% <= 0.303 milliseconds (cumulative count 13118)
29.817% <= 0.407 milliseconds (cumulative count 29817)
44.885% <= 0.503 milliseconds (cumulative count 44885)
58.751% <= 0.607 milliseconds (cumulative count 58751)
69.084% <= 0.703 milliseconds (cumulative count 69084)
77.536% <= 0.807 milliseconds (cumulative count 77536)
83.382% <= 0.903 milliseconds (cumulative count 83382)
87.871% <= 1.007 milliseconds (cumulative count 87871)
90.940% <= 1.103 milliseconds (cumulative count 90940)
93.293% <= 1.207 milliseconds (cumulative count 93293)
94.745% <= 1.303 milliseconds (cumulative count 94745)
95.954% <= 1.407 milliseconds (cumulative count 95954)
96.793% <= 1.503 milliseconds (cumulative count 96793)
97.509% <= 1.607 milliseconds (cumulative count 97509)
97.997% <= 1.703 milliseconds (cumulative count 97997)
98.391% <= 1.807 milliseconds (cumulative count 98391)
98.708% <= 1.903 milliseconds (cumulative count 98708)
98.976% <= 2.007 milliseconds (cumulative count 98976)
99.151% <= 2.103 milliseconds (cumulative count 99151)
99.892% <= 3.103 milliseconds (cumulative count 99892)
99.993% <= 4.103 milliseconds (cumulative count 99993)
99.998% <= 5.103 milliseconds (cumulative count 99998)
100.000% <= 6.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 75131.48 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.631     0.088     0.543     1.327     2.023     5.207
====== SADD ======
  100000 requests completed in 1.33 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.095 milliseconds (cumulative count 1)
50.000% <= 0.551 milliseconds (cumulative count 51014)
75.000% <= 0.775 milliseconds (cumulative count 75114)
87.500% <= 0.991 milliseconds (cumulative count 87518)
93.750% <= 1.215 milliseconds (cumulative count 93796)
96.875% <= 1.463 milliseconds (cumulative count 96875)
98.438% <= 1.743 milliseconds (cumulative count 98468)
99.219% <= 1.999 milliseconds (cumulative count 99229)
99.609% <= 2.263 milliseconds (cumulative count 99613)
99.805% <= 2.567 milliseconds (cumulative count 99805)
99.902% <= 2.879 milliseconds (cumulative count 99904)
99.951% <= 3.287 milliseconds (cumulative count 99953)
99.976% <= 3.511 milliseconds (cumulative count 99976)
99.988% <= 3.839 milliseconds (cumulative count 99988)
99.994% <= 4.047 milliseconds (cumulative count 99994)
99.997% <= 4.183 milliseconds (cumulative count 99997)
99.998% <= 4.295 milliseconds (cumulative count 99999)
99.999% <= 4.351 milliseconds (cumulative count 100000)
100.000% <= 4.351 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.002% <= 0.103 milliseconds (cumulative count 2)
2.018% <= 0.207 milliseconds (cumulative count 2018)
12.531% <= 0.303 milliseconds (cumulative count 12531)
29.092% <= 0.407 milliseconds (cumulative count 29092)
44.080% <= 0.503 milliseconds (cumulative count 44080)
58.297% <= 0.607 milliseconds (cumulative count 58297)
68.800% <= 0.703 milliseconds (cumulative count 68800)
77.564% <= 0.807 milliseconds (cumulative count 77564)
83.527% <= 0.903 milliseconds (cumulative count 83527)
88.154% <= 1.007 milliseconds (cumulative count 88154)
91.225% <= 1.103 milliseconds (cumulative count 91225)
93.640% <= 1.207 milliseconds (cumulative count 93640)
95.250% <= 1.303 milliseconds (cumulative count 95250)
96.421% <= 1.407 milliseconds (cumulative count 96421)
97.188% <= 1.503 milliseconds (cumulative count 97188)
97.855% <= 1.607 milliseconds (cumulative count 97855)
98.310% <= 1.703 milliseconds (cumulative count 98310)
98.709% <= 1.807 milliseconds (cumulative count 98709)
99.007% <= 1.903 milliseconds (cumulative count 99007)
99.243% <= 2.007 milliseconds (cumulative count 99243)
99.394% <= 2.103 milliseconds (cumulative count 99394)
99.930% <= 3.103 milliseconds (cumulative count 99930)
99.996% <= 4.103 milliseconds (cumulative count 99996)
100.000% <= 5.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 75244.55 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.628     0.088     0.551     1.287     1.903     4.351
====== HSET ======
  100000 requests completed in 1.33 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.103 milliseconds (cumulative count 1)
50.000% <= 0.551 milliseconds (cumulative count 50476)
75.000% <= 0.775 milliseconds (cumulative count 75360)
87.500% <= 0.991 milliseconds (cumulative count 87675)
93.750% <= 1.231 milliseconds (cumulative count 93823)
96.875% <= 1.479 milliseconds (cumulative count 96914)
98.438% <= 1.735 milliseconds (cumulative count 98465)
99.219% <= 2.015 milliseconds (cumulative count 99230)
99.609% <= 2.327 milliseconds (cumulative count 99619)
99.805% <= 2.495 milliseconds (cumulative count 99806)
99.902% <= 2.735 milliseconds (cumulative count 99906)
99.951% <= 2.911 milliseconds (cumulative count 99953)
99.976% <= 3.183 milliseconds (cumulative count 99977)
99.988% <= 3.319 milliseconds (cumulative count 99988)
99.994% <= 3.495 milliseconds (cumulative count 99994)
99.997% <= 3.623 milliseconds (cumulative count 99997)
99.998% <= 4.183 milliseconds (cumulative count 100000)
100.000% <= 4.183 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.001% <= 0.103 milliseconds (cumulative count 1)
1.823% <= 0.207 milliseconds (cumulative count 1823)
11.904% <= 0.303 milliseconds (cumulative count 11904)
28.308% <= 0.407 milliseconds (cumulative count 28308)
43.524% <= 0.503 milliseconds (cumulative count 43524)
57.821% <= 0.607 milliseconds (cumulative count 57821)
68.661% <= 0.703 milliseconds (cumulative count 68661)
77.771% <= 0.807 milliseconds (cumulative count 77771)
83.766% <= 0.903 milliseconds (cumulative count 83766)
88.277% <= 1.007 milliseconds (cumulative count 88277)
91.152% <= 1.103 milliseconds (cumulative count 91152)
93.410% <= 1.207 milliseconds (cumulative count 93410)
94.953% <= 1.303 milliseconds (cumulative count 94953)
96.238% <= 1.407 milliseconds (cumulative count 96238)
97.080% <= 1.503 milliseconds (cumulative count 97080)
97.854% <= 1.607 milliseconds (cumulative count 97854)
98.332% <= 1.703 milliseconds (cumulative count 98332)
98.720% <= 1.807 milliseconds (cumulative count 98720)
98.991% <= 1.903 milliseconds (cumulative count 98991)
99.213% <= 2.007 milliseconds (cumulative count 99213)
99.363% <= 2.103 milliseconds (cumulative count 99363)
99.968% <= 3.103 milliseconds (cumulative count 99968)
99.998% <= 4.103 milliseconds (cumulative count 99998)
100.000% <= 5.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 74962.52 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.631     0.096     0.551     1.311     1.911     4.183
====== ZADD ======
  100000 requests completed in 1.32 seconds
  50 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.095 milliseconds (cumulative count 1)
50.000% <= 0.543 milliseconds (cumulative count 50168)
75.000% <= 0.767 milliseconds (cumulative count 75577)
87.500% <= 0.975 milliseconds (cumulative count 87819)
93.750% <= 1.199 milliseconds (cumulative count 93831)
96.875% <= 1.447 milliseconds (cumulative count 96906)
98.438% <= 1.711 milliseconds (cumulative count 98461)
99.219% <= 1.999 milliseconds (cumulative count 99225)
99.609% <= 2.327 milliseconds (cumulative count 99611)
99.805% <= 2.711 milliseconds (cumulative count 99805)
99.902% <= 3.055 milliseconds (cumulative count 99903)
99.951% <= 3.287 milliseconds (cumulative count 99952)
99.976% <= 3.631 milliseconds (cumulative count 99976)
99.988% <= 3.799 milliseconds (cumulative count 99988)
99.994% <= 3.887 milliseconds (cumulative count 99994)
99.997% <= 4.119 milliseconds (cumulative count 99997)
99.998% <= 4.143 milliseconds (cumulative count 99999)
99.999% <= 4.247 milliseconds (cumulative count 100000)
100.000% <= 4.247 milliseconds (cumulative count 100000)

Cumulative distribution of latencies:
0.001% <= 0.103 milliseconds (cumulative count 1)
1.847% <= 0.207 milliseconds (cumulative count 1847)
12.127% <= 0.303 milliseconds (cumulative count 12127)
28.969% <= 0.407 milliseconds (cumulative count 28969)
44.401% <= 0.503 milliseconds (cumulative count 44401)
58.899% <= 0.607 milliseconds (cumulative count 58899)
69.919% <= 0.703 milliseconds (cumulative count 69919)
78.540% <= 0.807 milliseconds (cumulative count 78540)
84.569% <= 0.903 milliseconds (cumulative count 84569)
88.995% <= 1.007 milliseconds (cumulative count 88995)
91.773% <= 1.103 milliseconds (cumulative count 91773)
93.969% <= 1.207 milliseconds (cumulative count 93969)
95.452% <= 1.303 milliseconds (cumulative count 95452)
96.563% <= 1.407 milliseconds (cumulative count 96563)
97.313% <= 1.503 milliseconds (cumulative count 97313)
97.978% <= 1.607 milliseconds (cumulative count 97978)
98.432% <= 1.703 milliseconds (cumulative count 98432)
98.766% <= 1.807 milliseconds (cumulative count 98766)
99.044% <= 1.903 milliseconds (cumulative count 99044)
99.238% <= 2.007 milliseconds (cumulative count 99238)
99.385% <= 2.103 milliseconds (cumulative count 99385)
99.915% <= 3.103 milliseconds (cumulative count 99915)
99.996% <= 4.103 milliseconds (cumulative count 99996)
100.000% <= 5.103 milliseconds (cumulative count 100000)

Summary:
  throughput summary: 75757.57 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.623     0.088     0.543     1.271     1.887     4.247
```

## redis-go Benchmarks

| Benchmark       | Iterations | ns/op     |
|-----------------|------------|-----------|
| BenchmarkRedisSet | 47782      | 24615 ns/op |
| BenchmarkRedisGet | 49344      | 24318 ns/op |
| BenchmarkRedisIncr| 48984      | 24458 ns/op |
| BenchmarkRedisDel | 49016      | 23999 ns/op |


## redis-64-3.0.504 Benchmarks

https://github.com/microsoftarchive/redis/releases/download/win-3.0.504/Redis-x64-3.0.504.zip

| Benchmark       | Iterations | ns/op     |
|-----------------|------------|-----------|
| BenchmarkRedisSet | 50744      | 22828 ns/op |
| BenchmarkRedisGet | 52902      | 22471 ns/op |
| BenchmarkRedisIncr| 52924      | 22658 ns/op |
| BenchmarkRedisDel | 52605      | 22439 ns/op |

# 演示

<img src=".\img\demonstrate.png" style="zoom: 50%" alt="demonstrate"/>