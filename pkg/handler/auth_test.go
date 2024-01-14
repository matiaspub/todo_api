package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/matiaspub/todo-api/pkg/service"
	mock_service "github.com/matiaspub/todo-api/pkg/service/mocks"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           entity.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"test","username":"test","password":"qwerty"}`,
			inputUser: entity.User{
				Name:     "test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		}, {
			name:                "Empty Fields",
			inputBody:           `{"username":"test","password":"qwerty"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		}, {
			name:      "Service Failure",
			inputBody: `{"name":"test","username":"test","password":"qwerty"}`,
			inputUser: entity.User{
				Name:     "test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
