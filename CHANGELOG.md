Release Notes - fae - Version v0.0.2.rc
=======================================

### New Feature

    . kvdb sharding by cluster with replicas
    . kvdb servant and load test

### Improvement

    . universal resource pool management
    . warmup, connect backend on startup
    . rpc server thread pool
    . calculate concurrent conns for each remote php host, so that we know which host is busy, alarming
    . hot cache, local duplication

###  Bug

    . servant stats does not match engine stats
    . rps stats seems to have problem
    . call latency stats seems to have problem
    . fail to stop the engine

### Todo

    . mongodb/memcache graceful degrade with Circuit Breaker
    . memcache also support backend pools by name
    . pool pkg bug
    . LRUcache switch to better vitess implementation
    . use framed transport for better performance 
    . user service with auto local caching
    . proxy of servant pooling
    . optimize mc/mg pool, better failure handling
    . bitmap, replicated consitent hash
    . cpuprof/memprof
    . db flush in batch
    . IO throughput of servants
    . QoS
    . SLA of servants
    . better restart mechanism, socket pair?
    . simulate a low bandwidth, timeout env
    . parse/explain each db query to find bad query, e,g. no indexed query
    . THeaderTransport

Release Notes - fae - Version v0.0.1.alpha
==========================================

### New Feature

    . lcache(local cache shared among processes)
    . memcache backend
    . dlog(distributed log)
    . query server version info through REST api
    . auto discover peers like ES via multicast
    . app level exception added for memcache miss
    . golang servant client so that we can chain the call to distributed servers
    . unix domain socket rpc server transport
    . auto discovery of peers and form a cluster to serve requests
    . memcache consistent hash
    . global uniq id generator servant added
    . session based local var(like thread local) so that servant can maintain state across request calls

### Improvement

    . source code was greatly refactored and better organized
    . request context added, so we can do auditting
    . extend thrift to be able to trace request origin by extending TServerSocket
    . recyleable mongodb pooling 
    . profiler sampling rate feature added
    . config tcp nodelay to turn on/off Nagle
    . use memcache flags to auto (un)serialize php object, will not serialize if primitive type
    . an instance can disable some service so that we can customize deployment
    . use https://github.com/rcrowley/go-metrics as internal stats
    . control max outstanding sessions num

### Bug

    . fixed race condition



----
