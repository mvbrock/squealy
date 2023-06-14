# Squealy

A scalable streaming SQL query engine. Uses:

* [ZeroMQ's gossip protocol](http://czmq.zeromq.org/czmq4-0:zgossip) for clustering
* [Hyrise's SQL parser](https://github.com/hyrise/sql-parser) for SQL parsing
* [RocksDB](https://github.com/facebook/rocksdb) for data persistence
