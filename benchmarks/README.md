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
**Command**: go test -bench=. -cpu=2,4,8 -benchtime=60s -timeout=0

## MariaDB
```
BenchmarkMariaDBSingleInsertFixedData-2           233269            311992 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertFixedData-4           228982            312275 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertFixedData-8           226855            312530 ns/op             752 B/op         21 allocs/op
BenchmarkMariaDBSingleInsertRandomData-2          227834            312153 ns/op             816 B/op         23 allocs/op
BenchmarkMariaDBSingleInsertRandomData-4          227857            313258 ns/op             816 B/op         23 allocs/op
BenchmarkMariaDBSingleInsertRandomData-8          229183            312220 ns/op             816 B/op         23 allocs/op
```

## TCP
```
BenchmarkTCPClientSendSingleMessage-2           15554428              4678 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage-4           16660903              4306 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage-8           16607952              4318 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage2KB-2        13755962              5318 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage2KB-4        14499967              4955 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage2KB-8        13972306              5140 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage4KB-2        12424875              5837 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage4KB-4        13212968              5436 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage4KB-8        12871070              5610 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage8KB-2        11462521              6279 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage8KB-4        12012882              6013 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage8KB-8        12390741              5787 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage16KB-2       11903547              6086 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage16KB-4       11247807              6380 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage16KB-8       10999038              6334 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage32KB-2        6025442             11890 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage32KB-4        5746462             12581 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage32KB-8        5724616             12679 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage64KB-2        2994360             24246 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage64KB-4        2697742             28783 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage64KB-8        2548231             27723 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage128KB-2       1488265             48403 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage128KB-4       1329224             53235 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage128KB-8       1312656             52941 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage256KB-2        742324             97220 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage256KB-4        707037            104683 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage256KB-8        721300            103995 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage512KB-2        366891            193989 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage512KB-4        359095            204075 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage512KB-8        339342            203246 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1MB-2          194988            384171 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1MB-4          180643            408243 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1MB-8          180158            411231 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage10MB-2          18566           3895965 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage10MB-4          18394           3930107 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage10MB-8          18314           3930304 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage100MB-2          1887          38556426 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage100MB-4          1864          38979951 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage100MB-8          1866          39150366 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1GB-2             181         395581466 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1GB-4             180         399766696 ns/op               0 B/op          0 allocs/op
BenchmarkTCPClientSendSingleMessage1GB-8             178         400584389 ns/op               0 B/op          0 allocs/op
```

## HTTP
```
BenchmarkHTTPClientSendSingleMessage-2           1000000             66719 ns/op            6091 B/op         67 allocs/op
BenchmarkHTTPClientSendSingleMessage-4           1000000             63271 ns/op            6100 B/op         67 allocs/op
BenchmarkHTTPClientSendSingleMessage-8           1000000             63520 ns/op            6216 B/op         67 allocs/op
BenchmarkHTTPClientSendSingleMessage2KB-2         900810             77706 ns/op           10458 B/op         72 allocs/op
BenchmarkHTTPClientSendSingleMessage2KB-4        1000000             72630 ns/op           10481 B/op         72 allocs/op
BenchmarkHTTPClientSendSingleMessage2KB-8        1000000             73725 ns/op           10705 B/op         72 allocs/op
BenchmarkHTTPClientSendSingleMessage4KB-2         731761             99274 ns/op           17828 B/op         78 allocs/op
BenchmarkHTTPClientSendSingleMessage4KB-4         807621             88791 ns/op           17984 B/op         78 allocs/op
BenchmarkHTTPClientSendSingleMessage4KB-8         792387             90238 ns/op           18394 B/op         78 allocs/op
BenchmarkHTTPClientSendSingleMessage8KB-2         650334            111292 ns/op           44290 B/op         81 allocs/op
BenchmarkHTTPClientSendSingleMessage8KB-4         631486            106435 ns/op           44959 B/op         81 allocs/op
BenchmarkHTTPClientSendSingleMessage8KB-8         651879            110383 ns/op           46070 B/op         81 allocs/op
BenchmarkHTTPClientSendSingleMessage16KB-2        567175            126616 ns/op           81692 B/op         83 allocs/op
BenchmarkHTTPClientSendSingleMessage16KB-4        553728            130191 ns/op           83191 B/op         83 allocs/op
BenchmarkHTTPClientSendSingleMessage16KB-8        531114            136300 ns/op           85479 B/op         84 allocs/op
BenchmarkHTTPClientSendSingleMessage32KB-2        427830            170114 ns/op          192436 B/op         86 allocs/op
BenchmarkHTTPClientSendSingleMessage32KB-4        391776            183716 ns/op          197378 B/op         88 allocs/op
BenchmarkHTTPClientSendSingleMessage32KB-8        366886            196967 ns/op          201865 B/op         89 allocs/op
BenchmarkHTTPClientSendSingleMessage64KB-2        315596            230259 ns/op          323979 B/op         89 allocs/op
BenchmarkHTTPClientSendSingleMessage64KB-4        284298            253967 ns/op          332128 B/op         91 allocs/op
BenchmarkHTTPClientSendSingleMessage64KB-8        269028            267971 ns/op          337272 B/op         93 allocs/op
BenchmarkHTTPClientSendSingleMessage128KB-2       205194            346632 ns/op          554813 B/op         92 allocs/op
BenchmarkHTTPClientSendSingleMessage128KB-4       189961            378911 ns/op          566721 B/op         96 allocs/op
BenchmarkHTTPClientSendSingleMessage128KB-8       185986            382922 ns/op          572385 B/op         97 allocs/op
BenchmarkHTTPClientSendSingleMessage256KB-2       111212            659198 ns/op         1235081 B/op        100 allocs/op
BenchmarkHTTPClientSendSingleMessage256KB-4       102249            707682 ns/op         1246647 B/op        103 allocs/op
BenchmarkHTTPClientSendSingleMessage256KB-8       108058            719285 ns/op         1252100 B/op        105 allocs/op
BenchmarkHTTPClientSendSingleMessage512KB-2        53978           1213059 ns/op         2596347 B/op        111 allocs/op
BenchmarkHTTPClientSendSingleMessage512KB-4        49196           1386575 ns/op         2603641 B/op        112 allocs/op
BenchmarkHTTPClientSendSingleMessage512KB-8        47948           1464197 ns/op         2609279 B/op        113 allocs/op
BenchmarkHTTPClientSendSingleMessage1MB-2          21388           3758526 ns/op         5314825 B/op        118 allocs/op
BenchmarkHTTPClientSendSingleMessage1MB-4          17445           4099613 ns/op         5315669 B/op        118 allocs/op
BenchmarkHTTPClientSendSingleMessage1MB-8          19362           3452282 ns/op         5317715 B/op        118 allocs/op
BenchmarkHTTPClientSendSingleMessage10MB-2          2894          23264136 ns/op        52325260 B/op        129 allocs/op
BenchmarkHTTPClientSendSingleMessage10MB-4          3939          22653774 ns/op        52338492 B/op        130 allocs/op
BenchmarkHTTPClientSendSingleMessage10MB-8          3363          18553445 ns/op        52338670 B/op        130 allocs/op
BenchmarkHTTPClientSendSingleMessage100MB-2          554         139628378 ns/op        615297828 B/op       138 allocs/op
BenchmarkHTTPClientSendSingleMessage100MB-4          544         141763974 ns/op        615298993 B/op       139 allocs/op
BenchmarkHTTPClientSendSingleMessage100MB-8          488         148327561 ns/op        615301868 B/op       139 allocs/op
BenchmarkHTTPClientSendSingleMessage1GB-2             67         940118097 ns/op        5736587661 B/op      148 allocs/op
BenchmarkHTTPClientSendSingleMessage1GB-4             61        1117175741 ns/op        5736595158 B/op      149 allocs/op
BenchmarkHTTPClientSendSingleMessage1GB-8             79        1161890698 ns/op        5736598274 B/op      149 allocs/op
```
