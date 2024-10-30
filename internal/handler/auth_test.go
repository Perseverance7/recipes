package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/service"
	mock_service "github.com/Eagoker/recipes/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, user recipes.User)

	testTable := []struct {
		name string
		inputBody string
		inputUser recipes.User
		mockBehaviour mockBehaviour
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: `{"username": "Test", "password": "test"}`,
			inputUser: recipes.User{
				Username: "Test",
				Password: "test",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, user recipes.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"id":1}`,

		},
		{
			name: "Empty fields",
			inputBody: `{"password": "test"}`,
			mockBehaviour: func(s *mock_service.MockAuthorization, user recipes.User) {},
			expectedStatusCode: 400,
			expectedRequestBody: `{"message":"invalid input body"}`,

		},
		{
			name: "Service failure",
			inputBody: `{"username": "Test", "password": "test"}`,
			inputUser: recipes.User{
				Username: "Test",
				Password: "test",
			},
			mockBehaviour: func(s *mock_service.MockAuthorization, user recipes.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode: 500,
			expectedRequestBody: `{"message":"service failure"}`,

		},
	}

	for _, testCase := range testTable{
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))


			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())

		})
	}
}