package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ddProgerGo/task-kaspi/internal/handler"
	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) SavePerson(person models.Person) error {
	args := m.Called(person)
	return args.Error(0)
}

func (m *MockPersonService) GetPersonByIIN(iin string) (*models.Person, error) {
	args := m.Called(iin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Person), args.Error(1)
}

func (m *MockPersonService) GetPeopleByName(name string, page, limit int) ([]models.Person, error) {
	args := m.Called(name, page, limit)
	return args.Get(0).([]models.Person), args.Error(1)
}

func TestGetPersonByIIN(t *testing.T) {
	mockService := new(MockPersonService)

	validIIN := "020304550283"
	person := &models.Person{IIN: validIIN, Name: "John Doe", Phone: "1234567890"}

	mockService.On("GetPersonByIIN", validIIN).Return(person, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "iin", Value: validIIN})

	logger := logrus.New()
	h := handler.NewPersonHandler(mockService, logger)
	h.GetPersonByIIN(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Person
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, *person, response)
}

func TestSavePerson(t *testing.T) {
	mockService := new(MockPersonService)

	mockService.On("SavePerson", mock.Anything).Return(nil)

	body := `{"IIN": "020304550283", "Name": "Dulat Nurmeden", "Phone": "1234567890"}`
	req := httptest.NewRequest(http.MethodPost, "/save-person", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	logger := logrus.New()
	h := handler.NewPersonHandler(mockService, logger)
	h.SavePerson(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
