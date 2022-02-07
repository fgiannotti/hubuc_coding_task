package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	ErrorMsg   string `json:"error"`
}

type RegistrationController struct {
	logger      *zap.SugaredLogger
	users       UsersRepo
	encrpytions Encryptions
}

func (controller *RegistrationController) HandleRegister(c *gin.Context) {
	request, err := getRegisterRequestFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorResponse{http.StatusBadRequest, "Invalid body", err.Error()})
		return
	}

	_, err = controller.users.Get(request.Username)
	if err != nil {
		if err == UserNotFoundError(request.Username) {
			errResponse := ErrorResponse{http.StatusConflict, "Username already used", err.Error()}
			controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
			c.JSON(http.StatusConflict, errResponse)
			return
		}
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error getting user from db", err.Error()}
		controller.logger.Infow(errResponse.Message, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	encryptedPwd, err := controller.encrpytions.Encrypt(request.Password)
	if err != nil {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error encrypting pwd", err.Error()}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	newUsr := User{Name: request.Username, Email: request.Email, EncryptedPwd: encryptedPwd}
	err = controller.users.Save(newUsr)
	if err != nil {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error saving usr", err.Error()}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func getRegisterRequestFromBody(c *gin.Context) (registerRequest, error) {
	request := registerRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		return registerRequest{}, err
	}
	return request, nil
}
