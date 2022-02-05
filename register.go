package main

import (
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

func HandleRegister(c *gin.Context) {
	request, err := getRegisterRequestFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorResponse{http.StatusBadRequest, "Invalid body", err.Error()})
		return
	}
	//TODO: check if username is not used

	c.JSON(http.StatusOK, request)
}

func getRegisterRequestFromBody(c *gin.Context) (registerRequest, error) {
	request := registerRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		return registerRequest{}, err
	}
	return request, nil
}
