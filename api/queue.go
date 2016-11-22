package api

import (
	"strconv"

	"github.com/smira/aptly/queue"
	"github.com/gin-gonic/gin"
)

// GET /queue
func apiQueueList(c *gin.Context) {
	queue := context.Queue()
	c.JSON(200, queue.GetTasks())
}

// POST /queue/clear
func apiQueueClear(c *gin.Context) {
	queue := context.Queue()
	queue.Clear()
	c.JSON(200, gin.H{})
}

// GET /queue/wait
func apiQueueWait(c *gin.Context) {
	queue := context.Queue()
	queue.Wait()
	c.JSON(200, gin.H{})
}

// GET /queue/:id
func apiQueueTaskByID(c *gin.Context) {
	q := context.Queue()
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 0)
	if err != nil {
		c.Fail(500, err)
		return
	}

	var task queue.Task
	task, err = q.GetTaskByID(int(id))
	if err != nil {
		c.Fail(500, err)
		return
	}

	c.JSON(200, task)
}
