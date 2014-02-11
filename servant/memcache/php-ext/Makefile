all:
	cc -o hash crc32.c

test:all
	./hash
	go run crc32.go

clean:
	rm -f hash
