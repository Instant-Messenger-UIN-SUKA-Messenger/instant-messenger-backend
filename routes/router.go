package routes

import (
  "github.com/gin-gonic/gin"
  "instant-messenger-backend/rabbitmq"
)

func SetupRouter(r *gin.Engine) {
  r.POST("/messages", rabbitmq.PublishToDatabase)

  r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hello, world! Connect to Instant Messenger App!",
    })
  })

}