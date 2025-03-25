[![Go Report Card](https://goreportcard.com/badge/github.com/rah-0/benchmarks)](https://goreportcard.com/report/github.com/rah-0/hyperion)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Hyperion
Hyperion is a distributed, memory-first database written in Golang and for Golang projects. 

# Planning
- [ ] Benchmarking
  - [x] Protocols
  - [x] Serializers
  - [x] Compression
  - [x] Databases
    - [x] MariaDB
    - [x] PostgreSQL
    - [x] SQLite
  - [x] Cache
    - [x] Redis
    - [x] Memcached
  - [ ] GPU
    - [ ] nVidia CUDA
    - [ ] AMD ROCm 
- [ ] Implementation
  - [ ] Configuration
    - [x] Args
    - [ ] Hot Reload
      - [ ] HTTP Endpoint
  - [x] Storage
    - [x] Generator
      - [x] Versioning
      - [x] Migrations
    - [x] Memory
    - [x] Disk
  - [ ] Modes
    - [ ] Single node
    - [ ] Multi node
      - [ ] Sharding
      - [ ] Duplication
      - [ ] Sharding + Duplication
  - [ ] GPU 
    - [ ] Memory
    - [ ] Processing

---

## Benchmarking
In order to avoid bloating the repo with packages that will only be used for the sole purpose of benchmarking, another repo was created, you can find it [here](https://github.com/rah-0/benchmarks).  
The goal of benchmarking first before writing any code is to set some bars regarding expectations, establish some baselines and have some sanity checks.

### Protocols
**TCP** was picked over HTTP for performance reasons but also because:   
- netpoller is used, which is an abstraction over the operating system's I/O multiplexing mechanisms, including epoll on Linux systems
- persistent connections will be established between nodes

### Serializers
**gob** was picked over JSON or Protobuf because:
- JSON lacks efficient support for Golang's time.Time precision
- Protobuf requires constant mapping between internal and Protobuf structs, it also requires additional tooling and setup for schema management and code generation
- Gob uses Go's reflection to serialize native types directly without manual schema definitions

### Compression
If at some point compression is needed, **brotli** will be used but careful considerations need to be taken:
- storage: is the processing bill more expensive than increasing the drive size?
- over the wire: what's the biggest amount of data a slow network could transfer?

---

## Implementation

### Configuration
Since the config can be big enough, JSON is preferred for flexibility.  

- Command-line arguments: 
  - `pathConfig` config.json file path containing the cluster configuration
    - if no argument is passed, it will use the environment variable: `HyperionPathConfig`
  - `forceHost` can be used to load a specific configuration
    - if no argument is passed, it will use the environment variable: `HyperionForceHost`
- Hot Reload: required to modify the config without causing downtime
  - HTTP Endpoint: a POST request that will send the updated JSON config

See a configuration sample [here](https://github.com/rah-0/hyperion/blob/master/config/config.json)

**Details and considerations:**
- all nodes will have the same config
- the config file path passed on startup must be editable to allow the node to update its own config
- the HTTP Endpoint will be available for all nodes and the node that receives the new config is in charge of propagating the changes to the rest of the nodes
- when adding a new node, that node will propagate the updated config to the rest since adding a new node always implies updating config
- to remove a node, a POST request can be done, this request can be done to any of the nodes, including to the one that is to be removed
- on startup, each node will validate its own config with the rest, if there is a conflict, manual resolution is required
- node specific configuration is targeted by the Host.Name attribute

### Storage
How and where the information will be saved

#### Generator
Given [the defined](https://github.com/rah-0/hyperion/blob/master/entities/entities.go) structs (tables), a generator is needed to:
- avoid reflection
- facilitate queries
- handle versioning and migrations

##### Versioning
[The defined](https://github.com/rah-0/hyperion/blob/master/entities/entities.go) structs will have versions. Every time there is a change like:
- struct deleted
- struct field renamed
- struct field removed
- struct field added

a new version will be created, versioning is required for migrations

##### Migrations

Required in order to avoid downtime. At any point in time, structs might be updated or even completely removed. Migrations allow custom behavior to happen when upgrading or downgrading.
Upgrading allows current data to migrate to a future state while downgrade allows older nodes to communicate with nodes from different version. There are some instances where migrations
will not save you from having to restart every single node and cause a downtime. Let's see few cases:

This case illustrates a breaking change where it is required for all nodes to be upgraded, let's say there's this struct:
```GO
type Person struct {
	Name    string
	Surname string
}
```
and we want to update it to:
```GO
type Person struct {
	FullName    string
}
```
While the upgrade path is clear, simply concat Name + " " + Surname] the downgrade path is not possible since you cannot trust space " " to split the data.  
What if there is someone with compound name/surname and already contains a space? The chance of getting back the same previous results is unlikely.

The following case illustrates how a migration would be successful, let's say we want to migrate from:
```GO
type Event struct {
    Date  time.Time
}
```
to:
```GO
type Event struct {
    Date  int64
}
```
While int64 is not a direct representation of time.Time, you can still convert int64 to time.Time and viceversa as long as you keep in mind that you're losing precision. 

**Details and considerations:**
The generator will create the migration signature functions for you along with the tests 
but the correct implementation is your responsibility.

### GPU
A very long term plan is to allow the usage of GPU's to perform database operations such as:
- Relational Operations: projection, join, sorting, aggregation, grouping
- Compute: compression/decompression, encoding/decoding

---

# Sacrifices

Every distributed database has its drawbacks. This is what **Hyperion** sacrifices for **performance**:
- security: connections between nodes are not encrypted. If you need security you will have to manage it at a network level
- models: the entities (tables) have to [live inside](https://github.com/rah-0/hyperion/blob/master/entities/entities.go) the repo to make migrations easier
