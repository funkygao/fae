fae - Fun App Engine
====================
It's a middleware polyglot RPC engine for enterprise SOA infrastructure.

         ____      __      ____ 
        ( ___)    /__\    ( ___)
         )__)    /(__)\    )__) 
        (__)    (__)(__)  (____)

[![Build Status](https://travis-ci.org/funkygao/fae.png?branch=master)](https://travis-ci.org/funkygao/fae)
                               
### Why SOA?

*   Seperation of concerns
    - make a standard to hire developers
        - php frontend(more)
        - middleware backend(less)
*   Contract(IDL) based instead of language dependent servcie api
    - each level developers have a common sense
    - communicate by contract instead of direct call
*   Reuse common code as service and transparently reuse infrastructure
    - refuse copy & paste bug
*   Centralized best practice
*   Centralized monitoring, auditting and profiling
    - easy to find problems
*   Context free
*   Independently deployable/testable
    - vital code should be more robust
    - can't have too much vital code
*   lessen tcp 3/4 way handshake overhead
    - long conn pooling
    - meke better use of mem
*   Scale easily
    - frontend(php) and middleware scale dependently
    - middleware is in charge of performance while frontend is in charge of biz logic
*   Encapsulated 
    - logically decoupled and not shared its internal state
*   Polyglot development
    - some programs are not web based, e,g. batch, can be implemented as any language you like
*   Most large scale site use SOA as infrastructure

### Features

*   use multicast to auto discover fae peers for delegation
*   local LRU cache shared among processes
*   memcache servant
*   mongodb servant with transaction support
*   distributed logger servant

### Requirement

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

php.ini

    extension="thrift_protocol.so"
    extension="apc.so"

### Architecture


        +----------------+  +----------------+  +----------------+
        | php-fpm worker |  | php-fpm worker |  | php-fpm worker |
        +----------------+  +----------------+  +----------------+
            |                       |                       |
             -----------------------------------------------
                                    |                        
                                    | tcp/unix socket
                                    |                        
                            +---------------+
                            |  faed daemon  |
                            +---------------+
                            |  local cache  | 
                            +---------------+
                                    |                        
                                    | tcp long connection pool(keepalive)
                                    |                        
             -----------------------------------------------
            |                       |                       | hierarchy proxy
        +----------------+  +----------------+  +----------------+
        | mongodb servers|  |memcache servers|  |   faed proxy   |
        +----------------+  +----------------+  +----------------+

