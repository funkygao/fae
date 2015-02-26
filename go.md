golang for high performance
===========================

### Pointers usage

* In general, don’t share built-in type values with a pointer.
  - string, int, etc

* Slice of pointers are not cache friendly

*  If a value is passed to a value method, an on-stack copy can be used instead of allocating on the heap

### escape analysis

https://golang.org/doc/faq#stack_or_heap

    When possible, the Go compilers will allocate variables that are local to a function in that function's stack frame.
    If a variable escapes, it has to be allocated on the garbage-collected heap; otherwise, it’s safe to put it on the stack. This applies to new(T) allocations as well.
    Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.

    go build -gcflags '-m'
    src/cmd/gc/esc.c

### inlining

### compiler flags

    go build -gcflags 

    -S  see the assembly produced by compiling this package
    -l  disable inlining
    -m  show escape analysis and inlining
    -g  output the steps a the compiler is a taking at a very low level

### GC trace

    env GODEBUG=gctrace=1 godoc -http=:6060

### GC internal

    sysmon
      |
      |------------ forcegcperiod = 2m; // If we go two minutes without a garbage collection, force one to run
      |                                 // forcegchelper() -> sweepone() -> markroot() -> scanblock()
      |                                 // gc_m, gc mark
      |             scavengelimit = 5m; // Scavenger return the unused MSpan to OS
      |                                 // MHeap_Scavenge()
      |                                 // madvise(v, n, MADV_DONTNEED)
      |
    runtime_init
      |
    main_init
      |
    main_main


### Scheduler tracing

    GODEBUG=schedtrace=1000 ./daemon/faed/faed

### Tips

#### Avoid slice dynamic grow while append

    a := [100]string
    b := a[:0]
    b = append(b, "a")
    b = append(b, "b")
