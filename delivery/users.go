package delivery

import (
	"errors"
	"github.com/fgiannotti/hubuc_coding_task/core/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UsersController struct {
	logger      *zap.SugaredLogger
	users       services.UsersRepo
	encrpytions services.Encryptions
}

func NewUsersController(logger *zap.SugaredLogger, users services.UsersRepo, encrpytions services.Encryptions) UsersController {
	return UsersController{logger, users, encrpytions}
}

func (controller *UsersController) HandleGetUser(c *gin.Context) {
	username := c.Param("username")

	usr, err := controller.users.Get(username)
	if err != nil {
		if errors.Is(err,services.UserNotFoundError) {
			errResponse := ErrorResponse{http.StatusNotFound, "Username already used", err.Error()}
			controller.logger.Infow(errResponse.ErrorMsg, "username", username)
			c.JSON(http.StatusConflict, errResponse)
			return
		}
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error getting user from db", err.Error()}
		controller.logger.Infow(errResponse.Message, "username", username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}


	c.JSON(http.StatusCreated, usr)
}
