package delivery

import (
	"errors"
	"github.com/fgiannotti/hubuc_coding_task/core/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (controller *UsersController) HandleLogin(c *gin.Context) {
	request, err := getLoginRequestFromBody(c)
	if err != nil {
		errResponse := ErrorResponse{http.StatusBadRequest, "Invalid body", err.Error()}
		controller.logger.Info(errResponse.ErrorMsg)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	usr, err := controller.users.Get(request.Username)
	if err != nil {
		if errors.Is(err,services.UserNotFoundError) {
			errResponse := ErrorResponse{http.StatusNotFound, "User not found", err.Error()}
			controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
			c.JSON(errResponse.StatusCode, errResponse)
			return
		}
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error getting user from db", err.Error()}
		controller.logger.Infow(errResponse.Message, "username", request.Username)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	ok, err := controller.encrpytions.Compare(usr.EncryptedPwd, request.Password)
	if err != nil {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error checking pwd", err.Error()}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}
	if !ok {
		errResponse := ErrorResponse{http.StatusUnauthorized, "Error checking pwd", ""}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(errResponse.StatusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, "logged in")
}

func getLoginRequestFromBody(c *gin.Context) (loginRequest, error) {
	request := loginRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		return loginRequest{}, err
	}
	return request, nil
}
