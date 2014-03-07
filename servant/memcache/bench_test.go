package memcache

import (
	"testing"
)

func BenchmarkLegalKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		legalKey("errtag_asdfasdfasdfasfd_asdfaf")
	}
}

func BenchmarkStandardPickServer(b *testing.B) {
	client := getClient("standard", "127.0.0.1:11211", "127.0.0.1:11212", "127.0.0.1:11213")
	for i := 0; i < b.N; i++ {
		client.selector.PickServer("error_tag_232323232")
	}
}
