package gcar

// Item use for construct model to store cache
type Item struct {
	Object     interface{}
	Expiration int64
}

// Set should be set cache to memory
func Set() (string, bool) {
	return "Hello, World", true
}
