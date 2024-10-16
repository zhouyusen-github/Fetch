package main

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/add", addPoints)
	r.POST("/spend", spendPoints)
	r.GET("/balance", getPointsBalance)

	r.Run(":8000")
}

type AddRequest struct {
	Payer     string    `json:"payer"`
	Points    int       `json:"points"`
	Timestamp time.Time `json:"timestamp"`
}

type SpendRequest struct {
	Points int `json:"points"`
}

type SpendResult struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}

var (
	addRequests []AddRequest
	balance     = make(map[string]int)
	totalPoints int
	mu          sync.Mutex
)

func addPoints(c *gin.Context) {
	var addReq AddRequest
	if err := c.ShouldBindJSON(&addReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	addRequests = append(addRequests, addReq)
	balance[addReq.Payer] += addReq.Points
	totalPoints += addReq.Points

	c.Status(http.StatusOK)
}

func spendPoints(c *gin.Context) {
	var spendReq SpendRequest
	if err := c.ShouldBindJSON(&spendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if spendReq.Points > totalPoints {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn’t have enough points"})
		return
	}

	sort.Slice(addRequests, func(i, j int) bool {
		return addRequests[i].Timestamp.Before(addRequests[j].Timestamp)
	})

	var deleteBalance = make(map[string]int)

	i := 0
	for spendReq.Points > 0 {
		addReq := &addRequests[i]
		spend := 0
		if addReq.Points < 0 {
			spend = addReq.Points
		} else {
			spend = min(spendReq.Points, addReq.Points, balance[addReq.Payer])
		}

		addReq.Points -= spend
		spendReq.Points -= spend
		balance[addReq.Payer] -= spend
		totalPoints -= spend

		deleteBalance[addReq.Payer] -= spend

		if addReq.Points == 0 {
			addRequests = append(addRequests[:i], addRequests[i+1:]...)
		} else {
			i++
		}
	}

	var results []SpendResult
	for key, value := range deleteBalance {
		results = append(results, SpendResult{Payer: key, Points: value})
	}

	c.JSON(http.StatusOK, results)
}

func getPointsBalance(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, balance)
}

func min(a, b, c int) int {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}
