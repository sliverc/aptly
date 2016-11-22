package queue

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
	// FINISHED when task is finished
	FINISHED
)

// Task represents as task in a queue encapsulates process code
type Task struct {
	output  *TaskOutput
	process func(out *TaskOutput) error
	Name    string
	ID      int
	Err     error `json:",omitempty"`
	State   State
}

// TaskOutput represents a safe standard output of task
type TaskOutput struct {
	mu     *sync.Mutex
	output *bytes.Buffer
}

func (t *TaskOutput) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.String()
}

func (t *TaskOutput) Write(p []byte) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.Write(p)
}

// WriteString writes string to output
func (t *TaskOutput) WriteString(s string) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.output.WriteString(s)
}
