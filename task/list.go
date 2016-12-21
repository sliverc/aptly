package task

import (
	"bytes"
	"fmt"
	"sync"
	"errors"
)

// List is handling list of processes and makes sure
// only one process is executed at the time
type List struct {
	*sync.Mutex
	tasks []*Task
	wgTasks map[int]*sync.WaitGroup
	wg *sync.WaitGroup
	// resources currently used by running tasks
	usedResources *ResourcesSet
	idCounter int
}

// NewList creates empty task list
func NewList() *List {
	list := &List{
		Mutex: &sync.Mutex{},
		tasks: make([]*Task, 0),
		wgTasks: make(map[int]*sync.WaitGroup),
		wg: &sync.WaitGroup{},
		usedResources: NewResourcesSet(),
	}
	return list
}

// GetTasks gets complete list of tasks
func (list *List) GetTasks() []Task {
	var tasks []Task
	list.Lock()
	for _, task := range list.tasks {
		tasks = append(tasks, *task)
	}

	list.Unlock()
	return tasks
}

// GetTaskByID returns task with given id
func (list *List) GetTaskByID(ID int) (Task, error) {
	list.Lock()
	tasks := list.tasks
	list.Unlock()

	for _, task := range tasks {
		if task.ID == ID {
			return *task, nil
		}
	}

	return Task{}, fmt.Errorf("Could not find task with id %v", ID)
}

// GetTaskOutputByID returns standard output of task with given id
func (list *List) GetTaskOutputByID(ID int) (string, error) {
	task, err := list.GetTaskByID(ID)

	if err != nil {
		return "", err
	}

	return task.output.String(), nil
}


// RunTaskInBackground creates task and runs it in background. It won't be run and an error
// returned if there is a running tasks which is using any needed resources already.
func (list *List) RunTaskInBackground(name string, resources []string, process func(out *Output) error) (Task, error) {

	list.Lock()
	defer list.Unlock()

	if list.usedResources.ContainsAny(resources) {
		return Task{}, errors.New("Other running task already uses needed resources. Aborting...")
	}

	list.idCounter++
	list.usedResources.Add(resources)
	wgTask := &sync.WaitGroup{}
	task := &Task{
		output:  &Output{mu: &sync.Mutex{}, output: &bytes.Buffer{}},
		process: process,
		Name:    name,
		ID:      list.idCounter,
		State:   IDLE,
	}

	list.tasks = append(list.tasks, task)
	list.wgTasks[task.ID] = wgTask
	list.usedResources.Add(resources)

	list.wg.Add(1)
	wgTask.Add(1)

	go func() {

		list.Lock()
		{
			task.State = RUNNING
		}
		list.Unlock()

		err := process(task.output)

		list.Lock()
		{
			if err != nil {
				fmt.Fprintf(task.output, "Task failed with error: %v\n", err)
				task.State = FAILED
			} else {
				fmt.Fprintln(task.output, "Task succeeded")
				task.State = SUCCEEDED
			}

			list.usedResources.Remove(resources)

			wgTask.Done()
			list.wg.Done()
		}
		list.Unlock()
	}()

	return *task, nil
}

// Clear removes finished tasks from list
func (list *List) Clear() {
	list.Lock()

	var tasks []*Task
	for _, task := range list.tasks {
		if task.State == IDLE || task.State == RUNNING {
			tasks = append(tasks, task)
		}
	}
	list.tasks = tasks

	list.Unlock()
}

// Wait waits till all tasks are processed
func (list *List) Wait() {
	list.wg.Wait()
}

// WaitForTaskByID waits for task with given id to be processed
func (list *List) WaitForTaskByID(ID int) error {
	list.Lock()
	wgTask, ok := list.wgTasks[ID]
	list.Unlock()
	if !ok {
		return fmt.Errorf("Could not find task with id %v", ID)
	}

	wgTask.Wait()
	return nil
}
