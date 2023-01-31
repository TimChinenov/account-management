package main

import (
	"example/account-management/router"
)

func main() {
	router := router.SetupRouter()

	router.Run("localhost:8080")
}
