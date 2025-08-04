package interfaces

type Flushable interface {
	Flush() error
	Close() error
}
