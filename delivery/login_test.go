package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
	mock_services "github.com/fgiannotti/hubuc_coding_task/core/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginEndpointReturns200(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockLoginReq := loginRequest{Username: mockUsr.Name,Password: "pwd123"}
	mockLoginReqBytes, _ := json.Marshal(mockLoginReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	encryptionMock.EXPECT().Compare("encryptedPwd123","pwd123").Return(true,nil).Times(1)

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockUsr.Name)).Return(mockUsr,nil).Times(1)

	loggerMock := zaptest.NewLogger(t).Sugar()
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(mockLoginReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK,w.Code)
	assert.NotNil(t, w.Body)
}


func TestRegisterEndpointReturns500WhenEncrpytionCheckFails(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockLoginReq := loginRequest{Username: mockUsr.Name,Password: "pwd123"}
	mockLoginReqBytes, _ := json.Marshal(mockLoginReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	encryptionMock.EXPECT().Compare("encryptedPwd123","pwd123").Return(false,mockError).Times(1)

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockUsr.Name)).Return(mockUsr,nil).Times(1)

	loggerMock := zaptest.NewLogger(t).Sugar()
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(mockLoginReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError,w.Code)
}


func TestRegisterEndpointReturns500WhenUsersRepoFails(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockLoginReq := loginRequest{Username: mockUsr.Name,Password: "pwd123"}
	mockLoginReqBytes, _ := json.Marshal(mockLoginReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockLoginReq.Username)).Return(domain.User{},mockError).Times(1)

	controller := NewUsersController(zaptest.NewLogger(t).Sugar(), usersMock, mock_services.NewMockEncryptions(ctrl))

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(mockLoginReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError,w.Code)
}

func TestRegisterEndpointReturns400WhenBodyIsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	usersMock := mock_services.NewMockUsersRepo(ctrl)
	loggerMock := zaptest.NewLogger(t).Sugar()

	w := httptest.NewRecorder()
	c,_ := gin.CreateTestContext(w)
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)


	controller.HandleLogin(c)

	assert.Equal(t, http.StatusBadRequest,w.Code)
	assert.NotNil(t, w.Body)
}
