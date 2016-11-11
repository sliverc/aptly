package queue

// Processor is abstraction of logic to be enqueued
type Processor interface {
	Process() error
}

// FuncProcessor is a simple processor implementation using given func
// as processing unit
type FuncProcessor struct {
	fn func() error
}
