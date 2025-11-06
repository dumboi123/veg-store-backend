package rest_test

import (
	"net/http"
	"testing"
	"veg-store-backend/injection"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/restful/handler"
	"veg-store-backend/test/service"
	"veg-store-backend/test/unit/injection_test"

	"github.com/stretchr/testify/assert"
)

type UserHandler struct {
	*HandlerTest[*handler.UserHandler, *service.MockUserService]
}

func setupUserHandlerTest() *UserHandler {
	mockService := new(service.MockUserService)
	mockHandler := handler.NewUserHandler(mockService)
	engine := injection_test.MockUserRoutes(mockHandler)

	handlerTest := NewHandlerTest[*handler.UserHandler, *service.MockUserService](engine, mockHandler, mockService)
	return &UserHandler{
		HandlerTest: handlerTest,
	}
}

func (testHandler *UserHandler) TestHello_success(test *testing.T) {
	testHandler.MockService.On("Greeting").Return("Hello Ben")
	httpRecorder := testHandler.Get(test, AppURI("/user/hello"))
	assert.Equal(test, http.StatusOK, httpRecorder.Code)
	assert.Contains(test, httpRecorder.Body.String(), "Hello Ben")
	testHandler.MockService.AssertExpectations(test)
}

func (testHandler *UserHandler) TestDetails_withNotFoundID_fail(test *testing.T) {
	testHandler.MockService.On("FindById", "123").Return(nil, core.Error.NotFound.User)
	responseRecorder := testHandler.Get(test, AppURI("/user/details/123"))

	var response dto.HttpResponse[any]
	testHandler.DecodeResponse(test, responseRecorder, &response)
	assert.Equal(test, http.StatusNotFound, response.HttpStatus)
	assert.Equal(test, nil, response.Data)
	testHandler.MockService.AssertExpectations(test)
}

func TestUserHandler(test *testing.T) {
	injection.Inject("test")
	mockHandler := setupUserHandlerTest()
	test.Run("TestHello_success", mockHandler.TestHello_success)
	test.Run("TestDetails_withNotFoundID_fail", mockHandler.TestDetails_withNotFoundID_fail)
}
