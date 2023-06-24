package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"zog/delivery/handlers"
	"zog/delivery/models"
	"zog/domain/entity"
	infrastructure "zog/repository/infrastructure"
	repository "zog/repository/user"
	usecase "zog/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	_ "github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var router = gin.Default()
var handler *handlers.UserHandler
var userRepo *repository.UserRepository
var db *gorm.DB

func init() {
	db, _ = infrastructure.ConnectToDB()
	userRepo = repository.NewUserRepository(db)
	userUsecase := usecase.NewUser(userRepo)
	handler = &handlers.UserHandler{UserUsecase: userUsecase}
}

func TestSignup(t *testing.T) {

	router.POST("/signup", handler.Signup)

	payload := models.Signup{
		FirstName: "User",
		LastName:  "Sample",
		Email:     "usersample@gmail.com",
		Phone:     "9876543210",
		Password:  "pass@123",
	}

	jsonPayload, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("error in http request", err)
	}
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	// assert.Equal(t, payload, recorder.Body)

}

// func TestLoginWithPassword(t *testing.T) {

// 	router.POST("/loginwithpassword", handler.LoginWithPassword)

// 	form := url.Values{}
// 	form.Set("phone", "9876543210")
// 	form.Set("password", "pass@123")

// 	request, err := http.NewRequest("POST", "/loginwithpassword", strings.NewReader(form.Encode()))
// 	if err != nil {
// 		fmt.Println("error in http request", err)
// 	}
// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	recorder := httptest.NewRecorder()

// 	router.ServeHTTP(recorder, request)

// 	assert.Equal(t, http.StatusOK, recorder.Code)
// }

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetByID(id int) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestGetByID(t *testing.T) {
	repoMock := &UserRepositoryMock{}

	expectedUser := &entity.User{ID: 1, FirstName: "John"}
	repoMock.On("GetByID", 1).Return(expectedUser, nil)

}
