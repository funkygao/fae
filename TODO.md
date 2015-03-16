### TODO

*   [ ] stress test with different payload size
*   [ ] rpc cumCalls/cumSessions has bug
    - stress -x 1 -c 100 -n 100000 -logoff
*   [ ] thrift use slab allocator to read string
*   thrift oneway feature in golang
*   [ ] rate limit with golib.ratelimiter
    - plan to use plugin mechanism
*   [ ] bad performance related blocks
    - getSession()
    - logger
*   [ ] gotools
    - benchcmp
    - callgraph
*   [ ] make use of annotation to auot generate code skeleton
    - https://github.com/funkygao/goannotation
*   [X] shard lru cache to lower mutex race
*   [X] fae dashboard
*   [ ] stress loop in c1, c2 to test throughput under different concurrencies
*   [ ] more strict test on zookeeper failure
*   [ ] proxy pool, test on borrow
*   [X] make all db column not nullable
*   [ ] better request tracing
*   [ ] session.profiler should not be pointer, reduce GC overhead
*   [ ] mysql periodically ping to avoid being closed 
*   [X] optimize mysql query, iterate each row to transform to string/null
*   [X] engine record all err msg counter
*   [ ] when disk is full, fae will get stuck because of logging component
*   [ ] use jumbo frame to increase MTU 1500 -> 9000 to increase tcp throughput
*   [ ] log rotate size, only keep history for N days
*   [X] engine plugin
*   [X] use golib/signal SignalProcess instead of server.SignalProcess
*   [X] Context has too many strings, discard some of them
*   [X] change ctx.rid from string to int64, proxy servant rid generation mechanism
*   [X] start fae, then restart remote peer, then call ServantByKey, see what happens
*   [X] bloom filter 
*   [X] unified err logging so that external alarming system can get notified
*   [X] mysql prepare stmt caching
    - http://dev.mysql.com/doc/refman/5.1/en/query-cache-operation.html
    - CLIENT_NO_SCHEMA, don't allow database.table.column
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
