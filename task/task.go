package task

import (
	"bytes"
	"sync"
)

// State task is in
type State int

const (
	// IDLE when task is waiting
	IDLE State = iota
	// RUNNING when task is running
	RUNNING
	// SUCCEEDED when task is successfully finished
	SUCCEEDED
	// FAILED when task failed
	FAILED
)

// Task represents as task in a queue encapsulates process code
type Task struct {
	output  *Output
	process func(out *Output) error
	Name    string
	ID      int
	State   State
}

// Output represents a safe standard output of task
type Output struct {
	mu     *sync.Mutex
	output *bytes.Buffer
}

func (t *Output) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.String()
}

func (t *Output) Write(p []byte) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.Write(p)
}

// WriteString writes string to output
func (t *Output) WriteString(s string) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.WriteString(s)
}
