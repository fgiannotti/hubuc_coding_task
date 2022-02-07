package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
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
var mockError = errors.New("Its over 9000")

func TestRegisterEndpointReturns200(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockRegisterReq := registerRequest{Username: mockUsr.Name,Email: mockUsr.Email,Password: "pwd123"}
	mockRegisterReqBytes, _ := json.Marshal(mockRegisterReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	encryptionMock.EXPECT().Encrypt(mockRegisterReq.Password).Return("encryptedPwd123",nil).Times(1)

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockUsr.Name)).Return(domain.User{},nil).Times(1)
	usersMock.EXPECT().Save(gomock.Eq(mockUsr)).Return(nil).Times(1)

	loggerMock := zaptest.NewLogger(t).Sugar()
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(mockRegisterReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated,w.Code)
	assert.NotNil(t, w.Body)
}

func TestLoginEndpointReturns400WhenUsernameAlreadyExists(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockRegisterReq := registerRequest{Username: mockUsr.Name,Email: mockUsr.Email,Password: "pwd123"}
	mockRegisterReqBytes, _ := json.Marshal(mockRegisterReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockUsr.Name)).Return(mockUsr,nil).Times(1)

	controller := NewUsersController(zaptest.NewLogger(t).Sugar(), usersMock, mock_services.NewMockEncryptions(ctrl))

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(mockRegisterReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest,w.Code)
}

func TestLoginEndpointReturns500WhenEncryptionFails(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockRegisterReq := registerRequest{Username: mockUsr.Name,Email: mockUsr.Email,Password: "pwd123"}
	mockRegisterReqBytes, _ := json.Marshal(mockRegisterReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	encryptionMock.EXPECT().Encrypt(mockRegisterReq.Password).Return("",mockError).Times(1)

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockUsr.Name)).Return(domain.User{},nil).Times(1)

	loggerMock := zaptest.NewLogger(t).Sugar()
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(mockRegisterReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError,w.Code)
}


func TestLoginEndpointReturns500WhenUsersRepoFails(t *testing.T) {
	mockUsr := domain.User{"usr1","fgiannotti@pedidosya.com","encryptedPwd123"}
	mockRegisterReq := registerRequest{Username: mockUsr.Name,Email: mockUsr.Email,Password: "pwd123"}
	mockRegisterReqBytes, _ := json.Marshal(mockRegisterReq)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersMock := mock_services.NewMockUsersRepo(ctrl)
	usersMock.EXPECT().Get(gomock.Eq(mockRegisterReq.Username)).Return(domain.User{},mockError).Times(1)

	controller := NewUsersController(zaptest.NewLogger(t).Sugar(), usersMock, mock_services.NewMockEncryptions(ctrl))

	router := SetupRouter(controller)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(mockRegisterReqBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError,w.Code)
}

func TestLoginEndpointReturns400WhenBodyIsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptionMock := mock_services.NewMockEncryptions(ctrl)
	usersMock := mock_services.NewMockUsersRepo(ctrl)
	loggerMock := zaptest.NewLogger(t).Sugar()

	w := httptest.NewRecorder()
	c,_ := gin.CreateTestContext(w)
	controller := NewUsersController(loggerMock, usersMock, encryptionMock)


	controller.HandleRegister(c)

	assert.Equal(t, http.StatusBadRequest,w.Code)
	assert.NotNil(t, w.Body)
}
