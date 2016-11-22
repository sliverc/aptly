package queue

import (
	"bytes"
	"sync"
)

// Queue is handling list of processes and makes sure
// only one process is executed at the time
type Queue struct {
	*sync.Mutex
	work  chan *Task
	tasks []*Task
	// wait group to be able to close queue
	wgQueue sync.WaitGroup
	// wait group for tasks to finish
	wgtasks   sync.WaitGroup
	idCounter int
}

// New creates empty queue ready to be tasked
func New() *Queue {
	q := &Queue{
		Mutex: &sync.Mutex{},
		work:  make(chan *Task),
		tasks: make([]*Task, 0),
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

			q.Lock()
			task.State = RUNNING
			q.Unlock()

			err := task.process(task.output)

			q.Lock()
			task.Err = err
			task.State = FINISHED
			q.Unlock()

			q.wgtasks.Done()
		}
	}()

	return q
}

// GetTasks gets complete list of tasks
func (q *Queue) GetTasks() []*Task {
	q.Lock()
	tasks := q.tasks
	q.Unlock()
	return tasks
}

// Push pushes a new task with given name and processor logic to queue
func (q *Queue) Push(name string, process func(out *TaskOutput) error) {

	q.Lock()
	q.idCounter++
	task := &Task{
		output:  &TaskOutput{mu: &sync.Mutex{}, output: &bytes.Buffer{}},
		process: process,
		Name:    name,
		ID:      q.idCounter,
		State:   IDLE,
		Err:     nil,
	}
	q.tasks = append(q.tasks, task)
	q.Unlock()

	q.wgtasks.Add(1)
	go func() {
		q.work <- task
	}()
}

// Wait waits till all tasks are done on queue
func (q *Queue) Wait() {
	q.wgtasks.Wait()
}

// Clear removes finished tasks from list
func (q *Queue) Clear() {
	q.Lock()

	var tasks []*Task
	for _, task := range q.tasks {
		if task.State != FINISHED {
			tasks = append(tasks, task)
		}
	}
	q.tasks = tasks

	q.Unlock()
}

// Close stops running any additional tasks and waits still current tasks is finished
func (q *Queue) Close() {
	q.Lock()

	close(q.work)
	q.wgQueue.Wait()
	q.tasks = make([]*Task, 0)

	q.Unlock()
}
