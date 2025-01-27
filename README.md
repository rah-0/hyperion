# Hyperion
Hyperion is a distributed, memory-first database written in Golang and for Golang projects. 



## Planning
- [ ] Benchmarking
  - [x] Protocols
  - [x] Serializers
  - [x] Compression
  - [ ] Databases
    - [ ] MariaDB
    - [ ] PostgreSQL
  - [ ] Cache
    - [ ] Redis
    - [ ] Memcached
  - [ ] GPU
    - [ ] nVidia CUDA
    - [ ] AMD ROCm 
- [ ] Design
- [ ] Implementation
  - [ ] Memory
  - [ ] Storage
  - [ ] GPU 
    - [ ] Memory
    - [ ] Processing



### Benchmarking
In order to avoid bloating the repo with packages that will only be used for the sole purpose of benchmarking, another repo was created, you can find it [here](https://github.com/rah-0/benchmarks).  
The goal of benchmarking first before writing any code is to set some bars regarding expectations, establish some baselines and have some sanity checks.

#### Protocols
**TCP** was picked over HTTP for performance reasons but also because:   
- netpoller is used, which is an abstraction over the operating system's I/O multiplexing mechanisms, including epoll on Linux systems
- persistent connections will be established between nodes

#### Serializers
**gob** was picked over JSON or Protobuf because:
- JSON lacks efficient support for Golang's time.Time precision
- Protobuf requires constant mapping between internal and Protobuf structs, it also requires additional tooling and setup for schema management and code generation
- Gob uses Go's reflection to serialize native types directly without manual schema definitions

#### Compression
If at some point compression is needed, **brotli** will be used but careful considerations need to be taken:
- storage: is the processing bill more expensive than increasing the drive size?
- over the wire: what's the biggest amount of data a slow network could transfer?

### Implementation
#### GPU
A very long term plan is to allow the usage of GPU's to perform database operations such as:
- Relational Operations: projection, join, sorting, aggregation, grouping.
- Compute: compression/decompression, encoding/decoding
