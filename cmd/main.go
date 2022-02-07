package main

import (
	"github.com/fgiannotti/hubuc_coding_task/core/services"
	"github.com/fgiannotti/hubuc_coding_task/delivery"
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
	usersRepo := services.NewLocalUsersRepo(sugar)
	encryptionsService := services.NewBcryptEncryptionsService()

	usersController := delivery.NewUsersController(sugar, usersRepo, encryptionsService)

	r := delivery.SetupRouter(usersController)

	r.Run(":8080")
}