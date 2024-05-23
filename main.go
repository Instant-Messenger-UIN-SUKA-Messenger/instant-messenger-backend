package main

import (
	"github.com/gin-gonic/gin"
	"instant-messenger-backend/rabbitmq"
	"instant-messenger-backend/database"
	"instant-messenger-backend/routes"
)

func main() {
	r := gin.New()
	rabbitmq.ConnectToRabbitMQ()
	database.ConnectDB()
	routes.SetupRouter(r)
	r.Run("localhost:5000")
}
