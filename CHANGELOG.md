Release Notes - fae - Version v0.0.1.alpha

** New Feature

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
    . mongodb graceful degrade with Circuit Breaker
    . session based local var(like thread local)

** Improvement

    . source code was greatly refactored and better organized
    . request context added, so we can do auditting
    . extend thrift to be able to trace request origin by extending TServerSocket
    . recyleable mongodb pooling 
    . profiler sampling rate feature added
    . config tcp nodelay to turn on/off Nagle
    . use memcache flags to auto (un)serialize php object, will not serialize if primitive type
    . an instance can disable some service so that we can customize deployment
    . use https://github.com/rcrowley/go-metrics as internal stats

** Bug

    . fixed race condition

** Todo

    . use framed transport for better performance 
    . disable some services by config
    . user service with auto local caching
    . https://github.com/golang/glog
    . proxy of servant pooling
    . compress memcache data(check compress of live mc data)
    . optimize mc/mg pool
    . mongodb findColumn in php
    . bitmap, replicated consitent hash
    . cpuprof
    . memcache only as conn timeout, add io timeout

----
