package task

import (
	"errors"
	"testing"

	. "gopkg.in/check.v1"
)

// Launch gocheck tests
func Test(t *testing.T) {
	TestingT(t)
}

type QueueSuite struct{}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestQueue(c *C) {
	err := errors.New("Task failed")
	queue := NewQueue()
	c.Assert(len(queue.GetTasks()), Equals, 0)

	task := queue.Push("Successful task", func(out *Output) error {
		return nil
	})
	queue.Wait()

	tasks := queue.GetTasks()
	c.Assert(len(tasks), Equals, 1)
	task, _ = queue.GetTaskByID(task.ID)
	c.Check(task.State, Equals, SUCCEEDED)
	output, _ := queue.GetTaskOutputByID(task.ID)
	c.Check(output, Equals, "Task succeeded\n")

	task = queue.Push("Faulty task", func(out *Output) error {
		out.WriteString("Test Progress\n")
		return err
	})
	queue.Wait()

	tasks = queue.GetTasks()
	c.Assert(len(tasks), Equals, 2)
	task, _ = queue.GetTaskByID(task.ID)
	c.Check(task.State, Equals, FAILED)
	output, _ = queue.GetTaskOutputByID(task.ID)
	c.Check(output, Equals, "Test Progress\nTask failed with error: Task failed")
}
