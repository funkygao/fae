fproxy
======

funplus middleware layer which keeps long connection pool with backend and serves local php through unix domain socket

### Roles and Benefits

*   auditting for backend service
*   backend server location transparent for php(auto routing)
*   lesson tcp 3/4 way handshake overhead(conn pooling)
*   local cache(LRU)
*   profiler
*   queue for failed requests(auto retry)


### Architecture


        +----------------+  +----------------+  +----------------+
        | php-fpm worker |  | php-fpm worker |  | php-fpm worker |
        +----------------+  +----------------+  +----------------+
            |                       |                       |
             -----------------------------------------------
                                    |                        
                                    | unix domain socket
                                    |                        
                            +---------------+
                            | fproxy daemon |
                            +---------------+
                                    |                        
                                    | tcp long connection pool(keepalive)
                                    |                        
        +----------------+  +----------------+  +----------------+
        | mongodb servers|  |memcache servers|  | ... backends   |
        +----------------+  +----------------+  +----------------+

