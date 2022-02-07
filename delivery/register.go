package delivery

import (
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
	"github.com/fgiannotti/hubuc_coding_task/core/services"
	"github.com/gin-gonic/gin"
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

func (controller *UsersController) HandleRegister(c *gin.Context) {
	request, err := getRegisterRequestFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorResponse{http.StatusBadRequest, "Invalid body", err.Error()})
		return
	}

	usrFound, err := controller.users.Get(request.Username)
	if err != nil && err != services.UserNotFoundError(request.Username) {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error getting user from db", err.Error()}
		controller.logger.Infow(errResponse.Message, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	if usrFound.Name == request.Username {
		errResponse := ErrorResponse{http.StatusConflict, "Username already used", ""}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	encryptedPwd, err := controller.encrpytions.Encrypt(request.Password)
	if err != nil {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error encrypting pwd", err.Error()}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	newUsr := domain.User{Name: request.Username, Email: request.Email, EncryptedPwd: encryptedPwd}
	err = controller.users.Save(newUsr)
	if err != nil {
		errResponse := ErrorResponse{http.StatusInternalServerError, "Error saving usr", err.Error()}
		controller.logger.Infow(errResponse.ErrorMsg, "username", request.Username)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func getRegisterRequestFromBody(c *gin.Context) (registerRequest, error) {
	request := registerRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		return registerRequest{}, err
	}
	return request, nil
}
