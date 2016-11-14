package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smira/aptly/queue"
)

// GET /queue
func apiQueueList(c *gin.Context) {
	queue := context.Queue()
	c.JSON(200, queue.Tasks)
}

// POST /queue/clear
func apiQueueClear(c *gin.Context) {
	queue := context.Queue()
	queue.Clear()
	c.JSON(200, gin.H{})
}

// POST /queue/test
func apiQueueTest(c *gin.Context) {
	proc := queue.NewFuncProcessor(func() error {
		time.Sleep(1 * time.Minute)
		return nil
	})
	context.Queue().Push("/queue/test/", proc)

	c.JSON(202, gin.H{})
}
