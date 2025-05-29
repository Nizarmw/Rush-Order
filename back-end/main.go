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
	service.InitAdminSession(os.Getenv("SESSION_KEY"))
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500", "http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupSessionRoutes(r)
	routes.SetupCartRoutes(r)
	routes.SetupAdminRoutes(r, config.DB)
	routes.SetupProdukRoutes(r)
	routes.RegisterPaymentRoutes(r)

	r.Run(":8080")
}
