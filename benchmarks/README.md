# Benchmarks
These benchmarks are executed on a Virtual Machine.  
MariaDB is installed in the VM itself to avoid networking delays.

## About the system
CPU: AMD Ryzen 9 5950X  
RAM: 64GB (4x16GB) 4266MHz DDR4 CL19  

## About the VM
CPU: 8 threads  
RAM: 16GB  
OS: Debian 12, Kernel: 6.1.0-23-amd64  
NIC: 1GB  

## MariaDB
Version: 11.4.2-MariaDB-deb12  
Config:
```server.cnf
default-storage-engine=INNODB
sql-mode="STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION,NO_BACKSLASH_ESCAPES"

character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci
character-set-client-handshake=utf8mb4
default-time-zone=+00:00

join-buffer-size=512M
max-allowed-packet=2G
sort-buffer-size=128M
table-definition-cache=2000
max-connections=3000
tmp-table-size=512M

innodb-flush-log-at-trx-commit=1
innodb-log-buffer-size=512M
innodb-buffer-pool-size=16G
innodb-buffer-pool-instances=8
innodb-thread-concurrency=8
innodb-log-files-in-group=4
innodb-log-file-size=8G
innodb-write-io-threads=4
innodb-read-io-threads=4
innodb-autoextend-increment=256
innodb-old-blocks-time=500
innodb-file-per-table=ON

net-read-timeout=3600
net-write-timeout=3600
```

# Results

The aim of benching is to see what's faster while using the least amount of memory possible for each category.

**Test Command**: go test -bench=. -benchtime=60s -timeout=0

## DB's

---

## MariaDB
```
BenchmarkMariaDBSingleInsertFixedData-8           226855            312530 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8          229183            312220 ns/op             816 B/op         23 allocs/op
```

---

## Protocols (sending only)

| Message Size       | Protocol | Iterations   | Time per Op (ns) | Memory per Op (B)     | Allocations per Op |
|--------------------|----------|--------------|------------------|-----------------------|--------------------|
| **Single Message** |          |              |                  |                       |                    |
| Single Message-8   | TCP      | 16,607,952   | 4,318.0          | 0                     | 0                  |
| Single Message-8   | HTTP     | 1,000,000    | 63,520.0         | 6,216                 | 67                 |
| **2KB**            |          |              |                  |                       |                    |
| 2KB-8              | TCP      | 13,972,306   | 5,140.0          | 0                     | 0                  |
| 2KB-8              | HTTP     | 1,000,000    | 73,725.0         | 10,705                | 72                 |
| **4KB**            |          |              |                  |                       |                    |
| 4KB-8              | TCP      | 12,871,070   | 5,610.0          | 0                     | 0                  |
| 4KB-8              | HTTP     | 792,387      | 90,238.0         | 18,394                | 78                 |
| **8KB**            |          |              |                  |                       |                    |
| 8KB-8              | TCP      | 12,390,741   | 5,787.0          | 0                     | 0                  |
| 8KB-8              | HTTP     | 651,879      | 110,383.0        | 46,070                | 81                 |
| **16KB**           |          |              |                  |                       |                    |
| 16KB-8             | TCP      | 10,999,038   | 6,334.0          | 0                     | 0                  |
| 16KB-8             | HTTP     | 531,114      | 136,300.0        | 85,479                | 84                 |
| **32KB**           |          |              |                  |                       |                    |
| 32KB-8             | TCP      | 5,724,616    | 12,679.0         | 0                     | 0                  |
| 32KB-8             | HTTP     | 366,886      | 196,967.0        | 201,865               | 89                 |
| **64KB**           |          |              |                  |                       |                    |
| 64KB-8             | TCP      | 2,548,231    | 27,723.0         | 0                     | 0                  |
| 64KB-8             | HTTP     | 269,028      | 267,971.0        | 337,272               | 93                 |
| **128KB**          |          |              |                  |                       |                    |
| 128KB-8            | TCP      | 1,312,656    | 52,941.0         | 0                     | 0                  |
| 128KB-8            | HTTP     | 185,986      | 382,922.0        | 572,385               | 97                 |
| **256KB**          |          |              |                  |                       |                    |
| 256KB-8            | TCP      | 721,300      | 103,995.0        | 0                     | 0                  |
| 256KB-8            | HTTP     | 108,058      | 719,285.0        | 1,252,100             | 105                |
| **512KB**          |          |              |                  |                       |                    |
| 512KB-8            | TCP      | 339,342      | 203,246.0        | 0                     | 0                  |
| 512KB-8            | HTTP     | 47,948       | 1,464,197.0      | 2,609,279             | 113                |
| **1MB**            |          |              |                  |                       |                    |
| 1MB-8              | TCP      | 180,158      | 411,231.0        | 0                     | 0                  |
| 1MB-8              | HTTP     | 19,362       | 3,452,282.0      | 5,317,715             | 118                |
| **10MB**           |          |              |                  |                       |                    |
| 10MB-8             | TCP      | 18,314       | 3,930,304.0      | 0                     | 0                  |
| 10MB-8             | HTTP     | 3,363        | 18,553,445.0     | 52,338,670            | 130                |
| **100MB**          |          |              |                  |                       |                    |
| 100MB-8            | TCP      | 1,866        | 39,150,366.0     | 0                     | 0                  |
| 100MB-8            | HTTP     | 488          | 148,327,561.0    | 615,301,868           | 139                |
| **1GB**            |          |              |                  |                       |                    |
| 1GB-8              | TCP      | 178          | 400,584,389.0    | 0                     | 0                  |
| 1GB-8              | HTTP     | 79           | 1,161,890,698.0  | 5,736,598,274         | 149                |

## Encoding

| Size                 | Encoding Type | Iterations     | Time per Op (ns) | Memory per Op (B) | Allocations per Op |
|----------------------|---------------|----------------|------------------|-------------------|--------------------|
| **1Small-8**         |               |                |                  |                   |                    |
| 1Small-8             | GOB           | 331,346,038    | 218.2            | 24                | 1                  |
| 1Small-8             | JSON          | 302,926,762    | 235.1            | 24                | 1                  |
| 1Small-8             | Protobuf      | 675,272,739    | 106.9            | 80                | 1                  |
| **100Small-8**       |               |                |                  |                   |                    |
| 100Small-8           | GOB           | 12,789,297     | 5,647.0          | 24                | 1                  |
| 100Small-8           | JSON          | 5,796,138      | 12,506.0         | 24                | 1                  |
| 100Small-8           | Protobuf      | 6,710,169      | 10,709.0         | 8,000             | 100                |
| **10000Small-8**     |               |                |                  |                   |                    |
| 10000Small-8         | GOB           | 126,255        | 569,372.0        | 24                | 1                  |
| 10000Small-8         | JSON          | 56,889         | 1,268,709.0      | 97                | 1                  |
| 10000Small-8         | Protobuf      | 66,267         | 1,080,003.0      | 800,000           | 10,000             |
| **1000000Small-8**   |               |                |                  |                   |                    |
| 1000000Small-8       | GOB           | 1,198          | 60,117,976.0     | 24                | 1                  |
| 1000000Small-8       | JSON          | 561            | 128,433,139.0    | 26                | 1                  |
| 1000000Small-8       | Protobuf      | 658            | 110,364,613.0    | 80,000,002        | 1,000,000          |
| **1Unreal-8**        |               |                |                  |                   |                    |
| 1Unreal-8            | GOB           | 45,447,225     | 1,493.0          | 24                | 1                  |
| 1Unreal-8            | JSON          | 15,375,618     | 4,718.0          | 24                | 1                  |
| 1Unreal-8            | Protobuf      | 26,083,810     | 2,767.0          | 4,096             | 1                  |
| **10Unreals-8**      |               |                |                  |                   |                    |
| 10Unreals-8          | GOB           | 5,216,320      | 13,784.0         | 24                | 1                  |
| 10Unreals-8          | JSON          | 1,542,039      | 46,899.0         | 24                | 1                  |
| 10Unreals-8          | Protobuf      | 2,594,926      | 27,592.0         | 40,960            | 10                 |
| **100Unreals-8**     |               |                |                  |                   |                    |
| 100Unreals-8         | GOB           | 522,309        | 139,595.0        | 24                | 1                  |
| 100Unreals-8         | JSON          | 150,490        | 476,419.0        | 30                | 1                  |
| 100Unreals-8         | Protobuf      | 257,083        | 277,134.0        | 409,600           | 100                |
| **1000Unreals-8**    |               |                |                  |                   |                    |
| 1000Unreals-8        | GOB           | 48,220         | 1,513,438.0      | 24                | 1                  |
| 1000Unreals-8        | JSON          | 14,821         | 4,888,348.0      | 24                | 1                  |
| 1000Unreals-8        | Protobuf      | 25,994         | 2,790,997.0      | 4,096,000         | 1,000              |

### Observations:

- JSON: B/op discrepancies are most likely caused by GC.
- Protobuf: Memory and Allocations per op are quite high under heavy load, quotes from [Hacker News](https://news.ycombinator.com/item?id=40798740):

> Go gRPC server code does a lot of allocations. I have a gRPC service where each container does 50-80K/second of incoming calls and I spend a ton of time in GC and in allocating headers for all the msgs.



