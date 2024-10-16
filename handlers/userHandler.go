package handlers

import (
	"net/http"
	"sort"
	"sync"

	"Fetch/models"
	"Fetch/utils"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex

var (
	AddRequests []models.AddRequest
	Balance     = make(map[string]int)
	TotalPoints int
)

func AddPoints(c *gin.Context) {
	var addReq models.AddRequest
	if err := c.ShouldBindJSON(&addReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	AddRequests = append(AddRequests, addReq)
	Balance[addReq.Payer] += addReq.Points
	TotalPoints += addReq.Points

	c.Status(http.StatusOK)
}

func SpendPoints(c *gin.Context) {
	var spendReq models.SpendRequest
	if err := c.ShouldBindJSON(&spendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if spendReq.Points > TotalPoints {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn’t have enough points"})
		return
	}

	sort.Slice(AddRequests, func(i, j int) bool {
		return AddRequests[i].Timestamp.Before(AddRequests[j].Timestamp)
	})

	deleteBalance := make(map[string]int)
	i := 0
	for spendReq.Points > 0 {
		addReq := &AddRequests[i]
		var spend int
		if addReq.Points < 0 {
			spend = addReq.Points
		} else {
			spend = utils.Min(spendReq.Points, addReq.Points, Balance[addReq.Payer])
		}

		addReq.Points -= spend
		spendReq.Points -= spend
		Balance[addReq.Payer] -= spend
		TotalPoints -= spend

		deleteBalance[addReq.Payer] -= spend

		if addReq.Points == 0 {
			AddRequests = append(AddRequests[:i], AddRequests[i+1:]...)
		} else {
			i++
		}
	}

	var results []models.SpendResult
	for key, value := range deleteBalance {
		results = append(results, models.SpendResult{Payer: key, Points: value})
	}

	c.JSON(http.StatusOK, results)
}

func GetPointsBalance(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	c.JSON(http.StatusOK, Balance)
}
