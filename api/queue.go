package api

import (
	"github.com/gin-gonic/gin"
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

// GET /queue/wait
func apiQueueWait(c *gin.Context) {
	queue := context.Queue()
	queue.Wait()
	c.JSON(200, gin.H{})
}
