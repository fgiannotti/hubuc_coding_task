package delivery

import "github.com/gin-gonic/gin"

func SetupRouter(usersController UsersController) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", HealthCheck)
	r.POST("/users/register", usersController.HandleRegister)
	r.POST("/users/login", usersController.HandleLogin)
	r.GET("/users/:username", usersController.HandleGetUser)

	return r
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

