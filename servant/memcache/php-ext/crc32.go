package main

import "hash/crc32"

func main() {
    key := "test"
    hash := (crc32.ChecksumIEEE([]byte(key)) >> 16) & 0x7fff
    println(key, "->", hash)
}
