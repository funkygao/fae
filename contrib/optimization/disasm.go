package main

// go tool 6g -S disasm.go
// and check out the the assembler code to find the underhood
func main() {
	x := struct{ a, b int }{0x1, 0x2}
	b := x.a + x.b
	println(b)
}
