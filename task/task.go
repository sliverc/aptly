package task

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

// NewTask creates new task
func NewTask(process func(out *Output) error, name string, ID int) *Task {
	task := &Task{
		output:  NewOutput(),
		process: process,
		Name:    name,
		ID:      ID,
		State:   IDLE,
	}
	return task
}
