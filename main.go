package main

import (
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
	regController := RegistrationController{users: NewLocalUsersRepo(), logger: sugar, encrpytions: NewBcryptEncryptionsService()}
	r := gin.Default()
	r.GET("/ping", HealthCheck)
	r.POST("/users/register", regController.HandleRegister)
	r.POST("/users/login", HandleRegister)
	r.POST("/users/:username", HandleRegister)

	r.Run(":8080")
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
