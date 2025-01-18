package controller

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/christmas-fire/Bloomify/internal/service"
	mock_service "github.com/christmas-fire/Bloomify/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/golang/mock/gomock"
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
	// Init Test Table
	type mockBehavior func(r *mock_service.MockAuth, input signInInput)
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
				r.EXPECT().GenerateToken(input.Username, input.Password).Return("token", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"token":"token"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"username": "username"}`,
			inputSignIn:          signInInput{},
			mockBehavior:         func(r *mock_service.MockAuth, input signInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'signInInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "password": "qwerty"}`,
			inputSignIn: signInInput{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuth, input signInInput) {
				r.EXPECT().GenerateToken(input.Username, input.Password).Return("", errors.New("something went wrong"))
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
			test.mockBehavior(repo, test.inputSignIn)
			services := &service.Service{Auth: repo}
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
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
