package queue

import (
	"sync"
)

// TODO add unit test
// TODO rewrite api that writable calls are done through queue

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
	process func() error
	Name    string
	ID      int
	Err     error `json:",omitempty"`
	State   State
}

// Queue is handling list of processes and makes sure
// only one process is executed at the time
type Queue struct {
	mu        sync.Mutex
	work      chan *Task
	Tasks     []*Task
	// wait group to be able to close queue
	wgQueue   sync.WaitGroup
	// wait group for tasks to finish
	wgTasks   sync.WaitGroup
	idCounter int
}

// New creates empty queue ready to be tasked
func New() *Queue {
	q := &Queue{
		work:  make(chan *Task),
		Tasks: make([]*Task, 0),
	}

	// Start single worker for queue
	q.wgQueue.Add(1)
	go func() {
		for {
			task, ok := <-q.work
			if !ok {
				q.wgQueue.Done()
				return
			}

			q.mu.Lock()
			task.State = RUNNING
			q.mu.Unlock()

			err := task.process()

			q.mu.Lock()
			task.Err = err
			task.State = FINISHED
			q.mu.Unlock()

			q.wgTasks.Done()
		}
	}()

	return q
}

// Push pushes a new task with given name and processor logic to queue
func (q *Queue) Push(name string, process func() error) *Task {

	q.mu.Lock()
	q.idCounter++
	task := &Task{
		process: process, Name: name, ID: q.idCounter, State: IDLE, Err: nil,
	}
	q.Tasks = append(q.Tasks, task)
	q.mu.Unlock()

	q.wgTasks.Add(1)
	go func() {
		q.work <- task
	}()

	return task
}

// Wait waits till all tasks are done on queue
func (q *Queue) Wait() {
	q.wgTasks.Wait()
}

// Clear removes finished tasks from list
func (q *Queue) Clear() {
	q.mu.Lock()

	var tasks []*Task
	for _, task := range q.Tasks {
		if task.State != FINISHED {
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
	q.wgQueue.Wait()
	q.Tasks = make([]*Task, 0)

	q.mu.Unlock()
}
