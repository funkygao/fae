package main

import (
    "fmt"
    "labix.org/v2/mgo/bson"
)

var (
    in = []byte(`HuidZ/abtype$tutorialapaymentbgendarF`)
    out interface{}
)

func main() {
    if err := bson.Unmarshal(in, &out); err != nil {
        panic(err)
    }

    fmt.Printf("%#v\n", out)
}


