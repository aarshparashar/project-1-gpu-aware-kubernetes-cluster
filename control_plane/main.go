package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Models
type Job struct {
	ID        string `json:"id"`
	Priority  int    `json:"priority"`
	GPUReq    int    `json:"gpu_req"` // e.g., 1 GPU
	MemoryReq int    `json:"memory_req"`
	Status    string `json:"status"` // pending, running, completed
	NodeID    string `json:"node_id"`
}

type Node struct {
	ID          string    `json:"id"`
	GPUType     string    `json:"gpu_type"`
	GPUCapacity int       `json:"gpu_capacity"`
	GPUUsed     int       `json:"gpu_used"`
	LastSeen    time.Time `json:"last_seen"`
}

var (
	jobs  = make(map[string]*Job)
	nodes = make(map[string]*Node)
	mu    sync.RWMutex
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Job Endpoints
	r.POST("/api/v1/jobs", submitJob)
	r.GET("/api/v1/jobs/pending", getPendingJobs)
	r.PUT("/api/v1/jobs/:id/assign", assignJob)

	// Node Endpoints
	r.POST("/api/v1/nodes/register", registerNode)
	r.GET("/api/v1/nodes", getNodes)

	log.Println("Control Plane starting on port 8080...")
	r.Run(":8080")
}

func submitJob(c *gin.Context) {
	var job Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	job.Status = "pending"
	
	mu.Lock()
	jobs[job.ID] = &job
	mu.Unlock()

	c.JSON(http.StatusCreated, job)
}

func getPendingJobs(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	var pending []*Job
	for _, j := range jobs {
		if j.Status == "pending" {
			pending = append(pending, j)
		}
	}
	c.JSON(http.StatusOK, pending)
}

func assignJob(c *gin.Context) {
	jobID := c.Param("id")
	var req struct {
		NodeID string `json:"node_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	job, exists := jobs[jobID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	node, exists := nodes[req.NodeID]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node not found"})
		return
	}

	if node.GPUCapacity-node.GPUUsed < job.GPUReq {
		c.JSON(http.StatusConflict, gin.H{"error": "Node lacks GPU capacity"})
		return
	}

	job.Status = "running"
	job.NodeID = node.ID
	node.GPUUsed += job.GPUReq

	c.JSON(http.StatusOK, job)
}

func registerNode(c *gin.Context) {
	var node Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node.LastSeen = time.Now()
	
	mu.Lock()
	nodes[node.ID] = &node
	mu.Unlock()

	c.JSON(http.StatusOK, node)
}

func getNodes(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	var nodeList []*Node
	for _, n := range nodes {
		nodeList = append(nodeList, n)
	}
	c.JSON(http.StatusOK, nodeList)
}
