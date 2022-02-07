package main

import (
	"github.com/fgiannotti/hubuc_coding_task/core/services"
	"github.com/fgiannotti/hubuc_coding_task/delivery"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Running go server...")
	usersRepo := services.NewLocalUsersRepo()
	encryptionsService := services.NewBcryptEncryptionsService()

	usersController := delivery.NewUsersController(sugar, usersRepo, encryptionsService)

	r := gin.Default()
	r.GET("/ping", HealthCheck)
	r.POST("/users/register", usersController.HandleRegister)
	r.POST("/users/login", usersController.HandleLogin)
	r.POST("/users/:username", usersController.HandleGetUser)

	r.Run(":8080")
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
