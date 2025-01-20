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
**Command**: go test -bench=. -cpu=2,4,8 -benchtime=60s  
```
BenchmarkMariaDBSingleInsertFixedData-2           233269            311992 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertFixedData-4           228982            312275 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertFixedData-8           226855            312530 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-2          227834            312153 ns/op             816 B/op         23 allocs/op
BenchmarkMariaDBSingleInsertRandomData-4          227857            313258 ns/op             816 B/op         23 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8          229183            312220 ns/op             816 B/op         23 allocs/op

```
