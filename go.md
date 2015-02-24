golang for high performance
===========================

### escape analysis

https://golang.org/doc/faq#stack_or_heap

    When possible, the Go compilers will allocate variables that are local to a function in that function's stack frame.
    If a variable escapes, it has to be allocated on the garbage-collected heap; otherwise, it’s safe to put it on the stack. This applies to new(T) allocations as well.
    Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.

    go build -gcflags '-m'
    src/cmd/gc/esc.c
