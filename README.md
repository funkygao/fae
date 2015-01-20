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
                                    |                            peers discover
                                    |                            +---------------+
                                    |                     +------|  faed daemon  |-------+
                            +---------------+             |      +---------------+       |
                            |  faed daemon  |  tcp        |                              |
                            +---------------+ ------------|      peers discover          |
                            |  LRU cache    |  proxy      |      +---------------+       |
                            +---------------+             +------|  faed daemon  |-------|
                                    |                            +---------------+       |
                                    |                                                    |
                                    |                    consitent hash with replicas    |
                                    |----------------------------------------------------+
                                    |
                                    | tcp long connection pool(heartbeat) with recycling
                                    |
            +-----------------------------------------------+
            |                       |                       |     self contained
        +----------------+  +----------------+  +------------------------------+
        | mongodb servers|  |memcache servers|  | lcache | kvdb | idgen | ...  |
        +----------------+  +----------------+  +------------------------------+

### Why SOA?

*   More clear architecture
*   Seperation of concerns
*   Reuse common code as service and transparently reuse infrastructure
*   Centralized best practice
*   Centralized monitoring, auditting and profiling
*   Independently deployable/testable
    - vital code should be more robust
    - can't have too much vital code
*   Reduce tcp 3/4 way handshake overhead
*   Horizontal scale made easy
    - frontend(php) and middleware scale independently
    - middleware is in charge of performance while frontend is in charge of biz logic
*   Polyglot development

### Terms

*   Engine
*   Servant
*   Peer

### Highlights

*   Self manageable cluster
*   Dynamic cluster reconfiguration
    - VBucket
        - Better than consistent hashing
          - because they are easier to move between servers then individual keys
        - Never service a request on the wrong server
          - compared with consitent hash
        - Allow scaling up and down at will
        - We can hand data sets from one server another atomically
        - Servers still do not know about each other
*   Easy extending for more servants(RPC service)
*   Highly usage of mem to improve latancy & throughput
*   Circuit breaker protection
*   Full realtime internal stats export via http
*   Smart metrics with low overhead
*   Easy graceful degrade for OPS
    - auto
    - manual

### TODO

*   use of closed network connection
*   stress test on MacOS problems
    - Jan 19 12:43:51 mac-3 kernel[0]: process pingfae[84624] caught causing excessive wakeups. Observed wakeups rate (per sec): 2358; Maximum permitted wakeups rate (per sec): 150; Observation period: 300 seconds; Task lifetime number of wakeups: 45018
      - Mac sensors
      - sudo pmset -g
      - sudo pmset -a sms 0 # disable Sudden Motion Sensor
      - sudo pmset -a sms 1 # enable Sudden Motion Sensor
    - Limiting closed port RST response from 1422 to 250 packets per second
    - Jan 19 12:43:33 mac-3.local FileStatsAgent[84625]: Metadata.framework [Error]: couldn't get the client port
*   mysql prepare stmt caching
    - http://dev.mysql.com/doc/refman/5.1/en/query-cache-operation.html
    - CLIENT_NO_SCHEMA, don't allow database.table.column
    - too many round trips between fae and mysql
*   vBucket for cluster sharding, what about each kingdom is a shard?
*   hot reload config
*   more strict test on zookeeper failure
*   fae graceful shutdown
    - https://github.com/facebookgo/grace
*   maybe profiler sample rate is totally controlled by client
*   zk connection loss and session expiration
    - http://www.ngdata.com/so-you-want-to-be-a-zookeeper/
    - default zk session timeout: 2 * tickTime ~ 20 * tickTime
    - echo 'mntr' | nc localhost 2181
    - echo 'stat' | nc localhost 2181
*   race condition detector
*   NewTBufferedTransportFactory buffer size, and php config buf size
*   https://issues.apache.org/jira/browse/THRIFT-826 TSocket: Could not write
*   http://www.slideshare.net/renatko/couchbase-performance-benchmarking
*   golang uses /proc/sys/net/core/somaxconn as listener backlog
    - increase it if you need over 128(default) simultaneous outstanding connections
*   https://github.com/toddlipcon/gremlins

### Contribs

*   https://github.com/phunt/zktop
