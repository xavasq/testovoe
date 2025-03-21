package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"testovoe/internal/models"
	"testovoe/internal/service"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUserByID(ctx context.Context, id int64, user *models.User) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUserByID(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUserByID)
	r.PUT("/users/:id", h.UpdateUserByID)
	return r
}

func TestCreateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	user := models.User{Name: "test", Email: "test@example.com"}
	mockService.On("CreateUser", mock.Anything, &user).Return(nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"test"`)
	assert.Contains(t, w.Body.String(), `"email":"test@example.com"`)
	mockService.AssertExpectations(t)
}

func TestCreateUser_BadRequest(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"name": "test", "email":}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "некорректные данные")
	mockService.AssertNotCalled(t, "CreateUser")
}

func TestGetUserByID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	user := &models.User{ID: 1, Name: "test", Email: "test@example.com"}
	mockService.On("GetUserByID", mock.Anything, int64(1)).Return(user, nil)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
	assert.Contains(t, w.Body.String(), `"name":"test"`)
	mockService.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	mockService.On("GetUserByID", mock.Anything, int64(999)).Return((*models.User)(nil), service.ErrUserNotFound)

	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "пользователь не найден")
	mockService.AssertExpectations(t)
}

func TestUpdateUserByID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	user := models.User{Name: "updated", Email: "updated@example.com"}
	mockService.On("UpdateUserByID", mock.Anything, int64(1), &user).Return(nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "пользователь успешно обновлен")
	mockService.AssertExpectations(t)
}

func TestUpdateUserByID_BadID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("PUT", "/users/abc", bytes.NewBuffer([]byte(`{"name": "updated"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "неверный формат ID")
	mockService.AssertNotCalled(t, "UpdateUserByID")
}
