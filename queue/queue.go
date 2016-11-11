package queue

import (
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
	Processor
	Name  string
	ID    int
	Err   error
	State State
}

// Queue is handling list of processes and makes sure
// only one process is executed at the time
type Queue struct {
	mu    sync.Mutex
	work  chan *Task
	tasks []*Task
	wg    sync.WaitGroup
}

// Tasks gets list of all tasks (open, running and finished)
func (q *Queue) Tasks() []*Task {
	return q.tasks
}

// New creates empty queue ready to be tasked
func New() *Queue {
	q := &Queue{
		work: make(chan *Task),
		tasks: make([]*Task, 0),
	}

	// Start single worker of queue
	q.wg.Add(1)
	go func() {
		for {
			task, ok := <-q.work
			if !ok {
				q.wg.Done()
				return
			}

			q.mu.Lock()
			task.State = RUNNING
			q.mu.Unlock()

			err := task.Process()

			q.mu.Lock()
			task.Err = err
			task.State = FINISHED
			q.mu.Unlock()
		}
	}()

	return q
}

// Enqueue enqueues a new task with given name and processor logic
func (q *Queue) Enqueue(name string, proc Processor) *Task {

	q.mu.Lock()
	task := &Task{
		Processor: proc, Name: name, ID: len(q.tasks) + 1, State: IDLE,
	}
	q.tasks = append(q.tasks, task)
	q.mu.Unlock()

	go func() {
		q.work <- task
	}()

	return task
}

// Clear removes finished tasks from list
func (q *Queue) Clear() {
	q.mu.Lock()

	var tasks []*Task
	for _, task := range q.tasks {
		if task.State != FINISHED {
			tasks = append(tasks, task)
		}
	}
	q.tasks = tasks

	q.mu.Unlock()
}

// Close stops running any additional tasks and waits still current tasks is finished
func (q *Queue) Close() {
	q.mu.Lock()

	close(q.work)
	q.wg.Wait()
	q.tasks = make([]*Task, 0)

	q.mu.Unlock()
}
