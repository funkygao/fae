fae
===

         ____      __      ____ 
        ( ___)    /__\    ( ___)
         )__)    /(__)\    )__) 
        (__)    (__)(__)  (____)
                               
Funplus App Engine

It's middleware RPC engine.

### Why?

*   Seperation of concerns
*   Reuse common code as service and transparently reuse infrastructure
*   Centralized best practice
*   Centralized monitoring, auditting and profiling
*   lessen tcp 3/4 way handshake overhead(conn pooling)
*   Scale
*   Polyglot development

### Features

*   local LRU cache shared among processes
*   memcache servant
*   mongodb servant
*   distributed logger servant

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
                            |  faed daemon  |
                            +---------------+
                                    |                        
                                    | tcp long connection pool(keepalive)
                                    |                        
        +----------------+  +----------------+  +----------------+
        | mongodb servers|  |memcache servers|  | ... backends   |
        +----------------+  +----------------+  +----------------+

