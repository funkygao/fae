package memcache

var (
	crlf            = []byte("\r\n")
	space           = []byte(" ")
	resultStored    = []byte("STORED\r\n")
	resultNotStored = []byte("NOT_STORED\r\n")
	resultExists    = []byte("EXISTS\r\n")
	resultNotFound  = []byte("NOT_FOUND\r\n")
	resultDeleted   = []byte("DELETED\r\n")
	resultEnd       = []byte("END\r\n")

	resultClientErrorPrefix = []byte("CLIENT_ERROR ")
)

const (
	buffered = 8 // arbitrary buffered channel size, for readability

	StardardHashStrategy    = "standard"
	ConstistentHashStrategy = "consistent"
)
