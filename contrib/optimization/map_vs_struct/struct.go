package main

// go tool 6g -S struct.go
// and check out the the assembler code to find the underhood
func main() {
	x := struct{ a, b int }{0xaa, 0xbb}
	b := x.a + x.b
	println(b)
}
