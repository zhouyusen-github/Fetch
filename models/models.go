package models

import "time"

type AddRequest struct {
	Payer     string    `json:"payer" binding:"required"`
	Points    int       `json:"points" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

type SpendRequest struct {
	Points int `json:"points" binding:"required"`
}

type SpendResult struct {
	Payer  string `json:"payer" binding:"required"`
	Points int    `json:"points" binding:"required"`
}
