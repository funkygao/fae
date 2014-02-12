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

** Improvement

    . source code was greatly refactored and better organized
    . request context added, so we can do auditting
    . extend thrift to be able to trace request origin by extending TServerSocket
    . recyleable mongodb pooling 
    . profiler sampling rate feature added

** Bug

    . fixed race condition


** Todo

    . graceful degrade with Circuit Breaker
    . use framed transport for better performance 
    . user service with auto local caching
    . https://github.com/golang/glog
    . proxy of servant pooling
    . https://github.com/rcrowley/go-metrics

----
