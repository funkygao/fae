package servant

type Endpoint string

// Trace of a RPC call.
// distributed tracing learned from google Dapper.
type Call struct {
	Rid        int64 // request id
	Reason     string
	Id         int64 // Span id
	ParentId   int64
	Annotation []Annotation
	Debug      bool
}

type Annotation struct {
	Timestamp int64
	Value     string
	Host      Endpoint
	Duration  int32
}
