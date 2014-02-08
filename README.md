fae
===
Fun App Engine

It's a middleware RPC engine for enterprise SOA infrastructure.

         ____      __      ____ 
        ( ___)    /__\    ( ___)
         )__)    /(__)\    )__) 
        (__)    (__)(__)  (____)

[![Build Status](https://travis-ci.org/funkygao/fae.png?branch=master)](https://travis-ci.org/funkygao/fae)
                               
### Why SOA?

*   Seperation of concerns
*   Contract
*   Context free
*   Independently deployable
*   Independently testable
*   Reuse common code as service and transparently reuse infrastructure
*   Centralized best practice
*   Centralized monitoring, auditting and profiling
*   lessen tcp 3/4 way handshake overhead(conn pooling)
*   Scale easily
*   Encapsulated - logically decoupled and not shared its internal state
*   Polyglot development

### Features

*   use multicast to auto discover fae peers for request hash
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
            |                       |                       | hierarchy
        +----------------+  +----------------+  +----------------+
        | mongodb servers|  |memcache servers|  |   faed proxy   |
        +----------------+  +----------------+  +----------------+

