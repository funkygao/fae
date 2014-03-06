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
        | mongodb servers|  |memcache servers|  | lcache | kvdb | dlog | idgen |
        +----------------+  +----------------+  +------------------------------+

### Why SOA?

*   More clear architecture
    - UI -> php/python/... -> fae -> backend
    - instead of php/python -> backend
*   Seperation of concerns
    - make a standard to hire developers
        - each level developers have a common sense
        - php frontend(more)
        - middleware backend(less)
    - loose coupling            
*   Reuse common code as service and transparently reuse infrastructure
    - refuse copy & paste bug
*   Centralized best practice
*   Centralized monitoring, auditting and profiling
    - easy to find problems
*   Independently deployable/testable
    - vital code should be more robust
    - can't have too much vital code
*   Reduce tcp 3/4 way handshake overhead
    - long conn pooling
    - make better use of mem
*   Horizontal scale made easy
    - frontend(php) and middleware scale dependently
    - middleware is in charge of performance while frontend is in charge of biz logic
*   Polyglot development
    - some programs are not web based, e,g. batch, can be implemented as any language you like
*   Most large scale site use SOA as infrastructure

#### Terms and sub-directories

*   client
    - demonstration of how to call servants
    - stress test
    - currently supports python/go/php/java
*   config
    - configuration lib shared by engine and servant
*   daemon
    - the fae deamon
    - usually one instance per host
*   engine
    - load config file
    - export internal status through REST api
    - launch servants
    - implements thrift rpc server
    - launch thrift rpc serve
    - major version controller
*   etc
    - the config file 
*   http
    - Plugin based REST interface 
    - for monitor/control purpose
*   servant
    - RPC server side implementation
    - peer
        - other fae daemon that can be auto discovered through multicast
        - watchdog for health of peers
        - handles proxied requests from other(peer) fae
    -   proxy
        - stub of calling remote fae peers transparently
        - pooling

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
*   idmap servent...
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

    git.apache.org/thrift.git/lib/go/thrift

    git clone https://github.com/apache/thrift.git
    cd thrift
    ./bootstrap.sh
    ./configure --prefix=/opt/app/thrift
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

