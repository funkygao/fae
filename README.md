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
    - frontend(php) and middleware scale dependently
    - middleware is in charge of performance while frontend is in charge of biz logic
*   Polyglot development

#### Directories

*   client
    - multilingual example code of how to call fae
    - stress test
*   config
    - configuration lib shared by engine and servant
*   etc
    - the config file 
*   daemon
    - the fae deamon
*   engine
    - load config file
    - export internal status through REST api
    - launch servants
    - implements thrift rpc server
    - session based audit/profiler
    - handles timeout
    - launch thrift rpc serve
    - major version controller
*   http
    - Plugin based REST interface 
    - for monitor/control purpose
*   servant
    - manages circuit break
    - call based audit/profiler
    - RPC server side implementation
    - peer
        - other fae daemon that can be auto discovered through multicast
        - watchdog for health of peers
        - handles proxied requests from other(peer) fae
    -   proxy
        - stub of calling remote fae peers transparently
        - pooling

### Capacity Plan

#### Current stats

##### Request/Day

*   memcache pool
    - get   420,247,276 
    - set   172,524,762 
*   mongodb pool
    - insert     28,480,704   
    - query     594,480,553  
    - update    254,677,379  
    - delete     11,922,536 
    - getmore        68,673 
    - command   749,813,398
*   total
    - memcache 0.6 Billion
    - mongodb  0.9 Billion
    - total    1.5 Billion

##### Bandwidth

*   web
    - 10 times 80Mb/s = 800Mb/s

*   memcache pool
    - 2 times 20 = 40Mb/s

*   mongodb pool
    - 60 times 25Mb/s = 1.5Gb/s

#### Requirement for fae

If a single fae is deployed for the whole cluster, its capacity requirement:

*   qps
    - 20000 call/s

*   bandwidth
    - 800Mb/s

*   net conns
    - local tcp port used 7000
    - persistent tcp conns 2000

*   summary

                php
                 |
                 | 1000 concurrent conns
                 |
                fae
                 |
                 | pool size 20
                 | total 1500 persistent backend tcp conns
                 |
                 | total 6000 simultaneous memcache conns at most
                 | total 1000 simultaneous mongodb conns at most
                 |
           +----------------+
           |                |
        memcache#6      mongodb#60

### Highlights

*   Easy extending for more servants(RPC service)
*   Cluster based servants that can delegate(proxyed) to remote servants based on dup consitent hash
*   Use multicast to auto discover fae peers for delegation
*   Highly usage of mem to improve latancy & throughput
*   Merge recent requests to reduce backend service load
*   Circuit breaker protection
*   Fallback to mem when backend storage fails
    - requires session sticky to work
    - mem as response, and auto retry backend storage
    - when threshold retries reached, put to message queue for latter more retries
*   Easy graceful degrade for OPS
    - auto
    - manual

### Servants

*   idgen to generate global uniq id
*   local LRU cache shared among processes
*   memcache servant
*   mongodb servant with transaction support
*   distributed logger servant
*   idmap servant...
*   user account servant...

### Points of failure

*   rpc client app could crash
*   rpc client hardware could crash
*   rpc client network card could fail
*   network contention could cause timeouts
*   network elements such as routers could fail
*   transmission errors may lose messages
*   client and server versions may be incompatable
*   server network card could fail
*   server may have hardware problems
*   server software may crash
*   backend system such as database may become corrupted

### Remarks

*   golang uses /proc/sys/net/core/somaxconn as listener backlog
    - increase it if you need over 128(default) simultaneous outstanding connections

### Dependencies

hg

    sudo apt-get install mercurial

thrift above 0.9.0 which depends on flex

    sudo apt-get install flex

    go get git.apache.org/thrift.git/lib/go/thrift

    git clone -b 0.9.1 https://github.com/apache/thrift thrift-0.9.1
    cd thrift-0.9.1
    ./bootstrap.sh
    ./configure --prefix=/opt/app/thrift --with-cpp=no --with-erlang=no --with-c_glib=no --with-perl=no --with-ruby=no --with-haskell=no --with-d=no
    make
    make -k check
    sh test/test.sh
    make install

thrift_protocol.so

    cd thrift/lib/php/src/ext/
    phpize
    ./configure --with-php-config=/usr/local/php/bin/php-config
    make
    make test

php.ini

    extension="thrift_protocol.so"
    extension="apc.so"


### TODO

*   distribution of how many calls per session
*   session timeout seems not working
*   latency.call count > servant.calls, why?
*   maybe profiler sample rate is totally controlled by client
*   hot configuration reload
*   session timeout, what if php worker?
*   FunServantImpl.session(ctx).profiler.do
*   stats, e,g. in/out bytes, outstanding sessions, sessions by src ip,
*   realtime tracking of concurrent sessions by client host
*   session based profiler sampling, the whole session wholy sampled or not
*   if mongodb not existent in config file, and mg query arrives, fae dies
*   rate limit of connection of a given user
