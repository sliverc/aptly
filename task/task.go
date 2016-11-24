package task

import (
	"bytes"
	"sync"
	"fmt"
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
// which is compatbile to AptlyProgress
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

// Start is needed for progress compatability
func (t *Output) Start() {
	// Not implemented
}

// Shutdown is needed for progress compatability
func (t *Output) Shutdown() {
	// Not implemented
}

// Flush is needed for progress compatability
func (t *Output) Flush() {
	// Not implemented
}

// InitBar is needed for progress compatability
func (t *Output) InitBar(count int64, isBytes bool) {
	// Not implemented
}

// ShutdownBar is needed for progress compatability
func (t *Output) ShutdownBar() {
	// Not implemented
}

// AddBar is needed for progress compatability
func (t* Output) AddBar(count int) {
	// Not implemented
}

// SetBar sets current position for progress bar
func (t* Output) SetBar(count int) {
	// Not implemented
}

// Printf does printf in a safe manner
func (t* Output) Printf(msg string, a ...interface{}) {
	fmt.Fprintf(t, msg, a)
}

// ColoredPrintf does printf in a safe manner + newline
// currently are now colors supported.
func (t* Output) ColoredPrintf(msg string, a ...interface{}) {
	t.Printf(msg + "\n", a)
}
