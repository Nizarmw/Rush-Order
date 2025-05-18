package main

import (
	"RushOrder/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()


	r.Run(":8080")
}
