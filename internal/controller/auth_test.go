package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/service"
	mock_service "github.com/christmas-fire/Bloomify/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandler_signUp(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAuth, user models.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "email": "querty@mail", "password": "qwerty"}`,
			inputUser: models.User{
				Username: "username",
				Email:    "querty@mail",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuth, user models.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputUser:            models.User{},
			mockBehavior:         func(r *mock_service.MockAuth, user models.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "email": "querty@mail", "password": "qwerty"}`,
			inputUser: models.User{
				Username: "username",
				Email:    "querty@mail",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuth, user models.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockAuth(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Auth: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuth, input signInInput)

	// Создаем образец Tokens для тестов
	expectedTokens := service.Tokens{
		AccessToken:  "testAccessToken",
		RefreshToken: "testRefreshToken",
	}
	expectedTokensBody, _ := json.Marshal(expectedTokens)

	tests := []struct {
		name                 string
		inputBody            string
		inputSignIn          signInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputSignIn: signInInput{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuth, input signInInput) {
				// Ожидаем вызов GenerateToken с правильными аргументами и возвращаем Tokens
				r.EXPECT().GenerateToken(input.Username, input.Password).Return(expectedTokens, nil)
			},
			expectedStatusCode:   http.StatusOK,              // Статус 200
			expectedResponseBody: string(expectedTokensBody), // Ожидаем JSON с токенами
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputSignIn:          signInInput{},
			mockBehavior:         func(r *mock_service.MockAuth, input signInInput) {}, // No service call
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Key: 'signInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Invalid Credentials", // Тест на неверный логин/пароль
			inputBody: `{"username": "username", "password": "wrongpassword"}`,
			inputSignIn: signInInput{
				Username: "username",
				Password: "wrongpassword",
			},
			mockBehavior: func(r *mock_service.MockAuth, input signInInput) {
				// Ожидаем ошибку "invalid username or password" от сервиса
				r.EXPECT().GenerateToken(input.Username, input.Password).Return(service.Tokens{}, errors.New("invalid username or password"))
			},
			expectedStatusCode:   http.StatusUnauthorized, // Статус 401
			expectedResponseBody: `{"message":"invalid username or password"}`,
		},
		{
			name:      "Service Error", // Тест на другие ошибки сервиса
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputSignIn: signInInput{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuth, input signInInput) {
				// Ожидаем любую другую ошибку
				r.EXPECT().GenerateToken(input.Username, input.Password).Return(service.Tokens{}, errors.New("internal service error"))
			},
			expectedStatusCode:   http.StatusInternalServerError, // Статус 500
			expectedResponseBody: `{"message":"internal service error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_service.NewMockAuth(c)
			test.mockBehavior(repo, test.inputSignIn)
			services := &service.Service{Auth: repo} // Теперь должно работать
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/auth/sign-in", handler.signIn)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			require.Equal(t, test.expectedStatusCode, w.Code)
			require.JSONEq(t, test.expectedResponseBody, w.Body.String()) // Сравнение JSON
		})
	}
}

// Новый тест для refresh
func TestHandler_refresh(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuth, input refreshInput)

	expectedTokens := service.Tokens{
		AccessToken:  "newAccessToken",
		RefreshToken: "newRefreshToken",
	}
	expectedTokensBody, _ := json.Marshal(expectedTokens)

	tests := []struct {
		name                 string
		inputBody            string
		inputRefresh         refreshInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"refresh_token": "validRefreshToken"}`,
			inputRefresh: refreshInput{
				RefreshToken: "validRefreshToken",
			},
			mockBehavior: func(r *mock_service.MockAuth, input refreshInput) {
				r.EXPECT().RefreshToken(input.RefreshToken).Return(expectedTokens, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: string(expectedTokensBody),
		},
		{
			name:                 "Wrong Input (No Token)",
			inputBody:            `{}`, // Пустое тело
			inputRefresh:         refreshInput{},
			mockBehavior:         func(r *mock_service.MockAuth, input refreshInput) {}, // No service call
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid refresh token provided"}`,
		},
		{
			name:      "Invalid or Expired Token",
			inputBody: `{"refresh_token": "invalidOrExpiredToken"}`,
			inputRefresh: refreshInput{
				RefreshToken: "invalidOrExpiredToken",
			},
			mockBehavior: func(r *mock_service.MockAuth, input refreshInput) {
				// Может вернуть любую из этих ошибок
				r.EXPECT().RefreshToken(input.RefreshToken).Return(service.Tokens{}, errors.New("refresh session not found"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"refresh session not found"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"refresh_token": "validRefreshToken"}`,
			inputRefresh: refreshInput{
				RefreshToken: "validRefreshToken",
			},
			mockBehavior: func(r *mock_service.MockAuth, input refreshInput) {
				r.EXPECT().RefreshToken(input.RefreshToken).Return(service.Tokens{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_service.NewMockAuth(c)
			test.mockBehavior(repo, test.inputRefresh)
			services := &service.Service{Auth: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/auth/refresh", handler.refresh)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/refresh",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			require.Equal(t, test.expectedStatusCode, w.Code)
			require.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
