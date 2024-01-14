package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/matiaspub/todo-api/pkg/service"
	mock_service "github.com/matiaspub/todo-api/pkg/service/mocks"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, token string)

	tesstTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehaviour        mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "No Header",
			headerName:           "",
			mockBehaviour:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"Authorization header is required"}`,
		},
	}

	for _, testCase := range tesstTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(auth, testCase.token)

			handler := NewHandler(&service.Service{Authorization: auth})

			r := gin.New()
			r.POST("/protected", handler.userIdentity, func(ctx *gin.Context) {
				id, _ := ctx.Get(userCtx)
				ctx.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
