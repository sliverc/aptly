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

// NewFuncProcessor creates new FuncProcessor with given func
func NewFuncProcessor(fn func() error) *FuncProcessor {
	return &FuncProcessor{fn}
}

// Process wraps function
func (p *FuncProcessor) Process() error {
	return p.fn()
}
