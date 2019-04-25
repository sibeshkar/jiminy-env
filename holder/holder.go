package holder

// Holder holds a string value in a concurrency-safe manner.
type Holder interface {
	Get() string
	Set(string)
}
