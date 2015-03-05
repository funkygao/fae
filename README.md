fae - Fun App Engine
====================
Distributed middleware layer of multilingual RPC engine for enterprise SOA infrastructure.

         ____      __      ____ 
        ( ___)    /__\    ( ___)
         )__)    /(__)\    )__) 
        (__)    (__)(__)  (____)

[![Build Status](https://travis-ci.org/funkygao/fae.png?branch=master)](https://travis-ci.org/funkygao/fae)
                               
### Architecture


        +----------------+  +----------------+  +----------------+
        | php-fpm worker |  | php-fpm worker |  | php-fpm worker |
        +----------------+  +----------------+  +----------------+
            |                       |                       |
            +-----------------------------------------------+
                                    |                        
                                    | short lived tcp/unix socket                        
                                    |                        
                                    |                  
                                    |                            +---------------+
                                    |                     +------|  faed daemon  |-------+
                            +---------------+             |      +---------------+       |
                            |  faed daemon  |  tcp pool   |                              |
                            +---------------+ ------------|                              | peers
                            |  LRU cache    |  discovery  |      +---------------+       |
                            +---------------+             +------|  faed daemon  |-------|
                                    |                            +---------------+       |
                                    |                                                    |
                                    |                                         zookeeper  |
                                    |----------------------------------------------------+
                                    |
                                    | tcp conn pool
                                    |
            +-----------------------------------------------+
            |                       |                       |          SET model
        +----------------+  +----------------+  +------------------------------+
        | mongodb/mysql  |  | memcache/redis |  | lcache | kvdb | idgen | ...  |
        +----------------+  +----------------+  +------------------------------+

### Why SOA?

*   Seperation of concerns
*   Reuse common code as service and transparently reuse infrastructure
*   Centralized best practice, monitoring, auditting and profiling
*   Independently deployable/testable
    - vital code should be more robust
    - can't have too much vital code
*   Horizontal scale made easy
    - frontend(php) and middleware scale independently
    - middleware is in charge of performance while frontend is in charge of biz logic
*   Polyglot development
*   Easier dev team management

### Terms

*   Engine
    - handles system level issues
*   Servant
    - handles RPC services logic
*   Proxy
    - local stub of remote fae
*   Peer
    - a remote fae instance
*   Session
    - a RPC client tcp conn with fae
*   Call
    - a RPC call

### Highlights

*   Self manageable cluster
*   Linear scale by adding more fae instance
*   Dynamic cluster reconfiguration
    - VBucket
        - Better than consistent hashing
          - because they are easier to move between servers then individual keys
        - Never service a request on the wrong server
          - compared with consitent hash
        - Allow scaling up and down at will
        - We can hand data sets from one server another atomically
        - Servers still do not know about each other
*   Highly usage of mem to improve latancy & throughput
*   Circuit breaker protection
*   Flow control
*   Full realtime internal stats export via http
*   Smart metrics with low overhead
*   Graceful degrade for OPS
    - auto
    - manual

### Thrift Payload

    msgType = CALL | REPLY | EXCEPTION | ONEWAY

     0 1 2 3 4 5 6 7 8 9 a b c d e f  0 1 2 3 4 5 6 7 8 9 a b c d e f
    +----------------------------------------------------------------+
    |          version = 0x80010000 | msgType                        |
    +----------------------------------------------------------------+
    |          method name string len                                |
    +----------------------------------------------------------------+
    |          method name string itself ...                         |
    +----------------------------------------------------------------+
    |          seqId(int32)                                          |
    +----------------------------------------------------------------+


### TODO

*   [ ] use golib/signal SignalProcess instead of server.SignalProcess
*   [ ] replace config.engine.runWatchdog with server.WatchConfig
*   [ ] rename myslq pool to group
*   [ ] rate limit with golib.ratelimiter

*   [ ] engine pass tcpClient.RemoteAddr to Context, Servant will know the client better
*   [ ] proxy pool, test on borrow
*   [ ] rpc export call.all has bug
    - stress -x 1 -c 100 -n 100000 -logoff
*   [ ] name3 found dup names, bug
*   [ ] make all db column not nullable
*   [ ] better request tracing
*   [ ] backpressure
*   [ ] gm presence shows not only online, but also last sync time
*   [ ] session.profiler should not be pointer, reduce GC overhead
*   [X] optimize mysql query, iterate each row to transform to string/null
*   [X] engine record all err msg counter
*   [ ] when disk is full, fae will get stuck because of logging component
*   [ ] bad performance related blocks
    - getSession()
*   [ ] use jumbo frame to increase MTU 1500 -> 9000 to increase tcp throughput
*   [X] Context has too many strings, discard some of them
*   [ ] database/sql QueryRow for AR::get
*   [ ] go vet
*   [ ] log rotate size, only keep history for N days
*   [ ] periodically reload name3 from db
*   [ ] try not to use string as rpc func param, its costly to convert between []byte
*   [X] change ctx.rid from string to int64, proxy servant rid generation mechanism
*   [X] start fae, then restart remote peer, then call ServantByKey, see what happens
*   [X] bloom filter 
*   [X] unified err logging so that external alarming system can get notified
*   [ ] more strict test on zookeeper failure
*   [X] mysql prepare stmt caching
    - http://dev.mysql.com/doc/refman/5.1/en/query-cache-operation.html
    - CLIENT_NO_SCHEMA, don't allow database.table.column
    - too many round trips between fae and mysql
*   [ ] vBucket for cluster sharding, what about each kingdom is a shard?
*   [X] hot reload config
*   [X] fae graceful shutdown
    - https://github.com/facebookgo/grace
*   [X] maybe profiler sample rate is totally controlled by client
*   [X] zk connection loss and session expiration
    - http://www.ngdata.com/so-you-want-to-be-a-zookeeper/
    - default zk session timeout: 2 * tickTime ~ 20 * tickTime
    - echo 'mntr' | nc localhost 2181
    - echo 'stat' | nc localhost 2181
*   [X] NewTBufferedTransportFactory buffer size, and php config buf size
*   [X] golang uses /proc/sys/net/core/somaxconn as listener backlog
    - increase it if you need over 128(default) simultaneous outstanding connections
*   [X] thrift compiler didn't implement oneway in golang

#### Reference

*   aws ec2 packets-per-second (pps) maximum rate is 100k in+out
    - http://www.rightscale.com/blog/cloud-management-best-practices/benchmarking-load-balancers-cloud
*   RPS/RFS in linux
    - http://huoding.com/2013/10/30/296
*   https://github.com/phunt/zktop
*   https://github.com/toddlipcon/gremlins
*   network band width is cost problem while latency is physical constraint
