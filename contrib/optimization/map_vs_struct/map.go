package main

// go tool 6g -S map.go
// and check out the the assembler code to find the underhood
func main() {
	x := map[string]int{"a": 0xaa, "b": 0xbb}
	b := x["a"] + x["b"]
	println(b)
}
