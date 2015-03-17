fae - Fun App Engine [![Build Status](https://travis-ci.org/funkygao/fae.png?branch=master)](https://travis-ci.org/funkygao/fae)
====================
Distributed RPC framework for enterprise SOA infrastructure.

         ____      __      ____ 
        ( ___)    /__\    ( ___)
         )__)    /(__)\    )__) 
        (__)    (__)(__)  (____)


**Table of Contents**

- [Usage](#usage)
- [Dashboard](#dashboard)
- [Architecture](#status)
- [Terms](#terms)
- [Highlights](#highlights)
- [Performance](#perf)
- [Reference](#reference)

### Usage

#### dependency

    install thrift 

#### compile

    ./build.sh

#### run

    # create a config file
    cp etc/etc/faed.cf.sample contrib
    ./contrib/build_cf.php # create the config file

    # startup fae
    ./daemon/faed/faed -conf contrib/faed.cf.rc
                               
### Dashboard

![dashboard](https://github.com/funkygao/fae/blob/master/contrib/resources/dashboard.png)

### Architecture


        +----------------+  +----------------+  +----------------+
        |   php-fpm      |  |    php-fpm     |  |     php-fpm    |
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
                            |  LRU cache    |  proxy      |      +---------------+       |
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

### Terms

*   Engine
    - handles system level issues
*   Servant
    - handles RPC services logic
*   Proxy
    - local stub of remote fae peer
*   Peer
    - a remote fae instance
*   Session
    - a RPC client tcp connection with fae
*   Call
    - a RPC call within a Session

### Highlights

*   Self manageable cluster with browser base dashboard
*   Linear scales by adding more fae instances
*   Highly usage of mem to improve latancy & throughput
*   Circuit breaker protection and flow control
*   Smart metrics with low overhead
*   Graceful degrade for OPS
*   Plugins
*   One binary, homogeneous deployment
*   Dynamic cluster reconfiguration with vbucket

### Performance

*   currently, a single fae node qps around 50k(no batch request)
    - limited by NIC PPS(packets per second)
    - has to write linux kernal module to overcome this
*   will be tweaked to 100k

#### Reference

*   aws ec2 packets-per-second (pps) maximum rate is 100k in+out
    - http://www.rightscale.com/blog/cloud-management-best-practices/benchmarking-load-balancers-cloud
*   RPS/RFS in linux
    - http://huoding.com/2013/10/30/296
*   http://highscalability.com/blog/2013/5/13/the-secret-to-10-million-concurrent-connections-the-kernel-i.html
*   https://svn.ntop.org/svn/ntop/trunk/PF_RING/
*   tcpcopy
*   golang
    - GODEBUG=schedtrace=1000
    - GODEBUG=gctrace=1

