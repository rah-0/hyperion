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
  - [x] Cache
    - [x] Redis
    - [x] Memcached
  - [ ] GPU
    - [ ] nVidia CUDA
    - [ ] AMD ROCm 
- [ ] Design
  - [x] Configuration
    - [x] Args
    - [x] Hot Reload
      - [x] HTTP Endpoint
- [ ] Implementation
  - [ ] Configuration
    - [ ] Args
    - [ ] Hot Reload
      - [ ] HTTP Endpoint
  - [ ] Memory
  - [ ] Storage
  - [ ] GPU 
    - [ ] Memory
    - [ ] Processing

---

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

---

### Design

#### Configuration
Since the config can be big enough, JSON is preferred for flexibility.  

- Startup Args: 
  - config.json file path containing the cluster configuration
  - data path where items will be stored
- Hot Reload: required to modify the config without causing downtime
  - HTTP Endpoint: a POST request that will send the updated JSON config

**Details and considerations:**
- all nodes will have the same config
- the config file path passed on startup must be editable to allow the node to update its own config
- the HTTP Endpoint will be available for all nodes and the node that receives the new config is in charge of propagating the changes to the rest of the nodes
- when adding a new node, that node will propagate the updated config to the rest since adding a new node always implies updating config
- to remove a node, a POST request can be done, this request can be done to any of the nodes, including to the one that is to be removed
- on startup, each node will validate its own config with the rest, if there is a conflict, manual resolution is required 

---

### Implementation

#### GPU
A very long term plan is to allow the usage of GPU's to perform database operations such as:
- Relational Operations: projection, join, sorting, aggregation, grouping.
- Compute: compression/decompression, encoding/decoding
