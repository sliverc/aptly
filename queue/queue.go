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
	name  string
	id    int
	err   error
	state State
}

// Queue is handling list of processes and makes sure
// only one process is executed at the time
type Queue struct {
	mu    sync.Mutex
	work  chan *Task
	Tasks []*Task
	wg    sync.WaitGroup
}

// New creates empty queue ready to be tasked
func New() *Queue {
	q := &Queue{
		work: make(chan *Task),
	}

	// Start single worker of queue
	q.wg.Add(1)
	go func() {
		task, ok := <-q.work
		if !ok {
			q.wg.Done()
			return
		}

		q.mu.Lock()
		task.state = RUNNING
		q.mu.Unlock()

		err := task.Process()

		q.mu.Lock()
		task.err = err
		task.state = FINISHED
		q.mu.Unlock()
	}()

	return q
}

// Enqueue enqueues a new task with given name and processor logic
func (q *Queue) Enqueue(name string, proc Processor) {

	q.mu.Lock()
	task := &Task{
		Processor: proc, name: name, id: len(q.Tasks) + 1, state: IDLE,
	}
	q.Tasks = append(q.Tasks, task)
	q.mu.Unlock()

	go func() {
		q.work <- task
	}()
}

// Clear removes finished tasks from list
func (q *Queue) Clear() {
	q.mu.Lock()

	var tasks []*Task
	for _, task := range q.Tasks {
		if task.state != FINISHED {
			tasks = append(tasks, task)
		}
	}
	q.Tasks = tasks

	q.mu.Unlock()
}

// Close stops running any additional tasks and waits still current tasks is finished
func (q *Queue) Close() {
	q.mu.Lock()

	close(q.work)
	q.wg.Wait()
	q.Tasks = make([]*Task, 0)

	q.mu.Unlock()
}
