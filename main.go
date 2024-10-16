package main

import (
	"Fetch/router"
)

func main() {
	router := router.InitRouter()
	router.Run(":8000")
}
