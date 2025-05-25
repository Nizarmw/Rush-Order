package main

import (
	"RushOrder/config"
	"RushOrder/routes"
	"RushOrder/service"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	log.Println("Starting RushOrder server...")
	config.InitDB()

	service.InitSessionStore(os.Getenv("SESSION_KEY"))

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupSessionRoutes(r)
	routes.SetupCartRoutes(r)
	r.Run(":8080")
}
