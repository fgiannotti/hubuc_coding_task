package delivery

import (
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
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
	request, err := getRegisterRequestFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorResponse{http.StatusBadRequest, "Invalid body", err.Error()})
		return
	}

	_, err = controller.users.Get(request.Username)
	if err != nil {
		if err == services.UserNotFoundError(request.Username) {
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
