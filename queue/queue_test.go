package queue

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
	queue := New()
	c.Assert(len(queue.GetTasks()), Equals, 0)

	queue.Push("Successful task", func(out *TaskOutput) error {
		return nil
	})
	queue.Wait()

	tasks := queue.GetTasks()
	c.Assert(len(tasks), Equals, 1)
	c.Check(tasks[0].State, Equals, FINISHED)

	queue.Push("Faulty task", func(out *TaskOutput) error {
		out.WriteString("Test Progress")
		return err
	})
	queue.Wait()

	tasks = queue.GetTasks()
	c.Assert(len(tasks), Equals, 2)
	c.Check(tasks[1].State, Equals, FINISHED)
	c.Check(tasks[1].Err, Equals, err)
	c.Check(tasks[1].output.String(), Equals, "Test Progress")
}
