package handler

import (
	"net/http"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/service"

	"go.uber.org/fx"
)

type AuthHandler struct {
	service service.AuthenticationService
}

func NewAuthHandler(authenticationService service.AuthenticationService) *AuthHandler {
	return &AuthHandler{service: authenticationService}
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.SignInRequest true "User credentials"
// @Success 200 {object} dto.HttpResponse[dto.Tokens]
// @Failure 401 {object} dto.HttpResponse[string]
// @Router /user/sign-in [post]
func (handler *AuthHandler) SignIn(context *core.HttpContext) {
	var request dto.SignInRequest
	var err error

	err = context.Gin.ShouldBindJSON(&request)
	tokens, err := handler.service.Tokens(request)
	if err != nil {
		context.Gin.Error(core.Error.Auth.Unauthenticated)
		return
	}

	context.JSON(http.StatusOK, dto.HttpResponse[dto.Tokens]{
		HttpStatus: http.StatusOK,
		Data:       *tokens,
	})
}

//// Info godoc
//// @Summary User details
//// @Description Get details of a user by id
//// @Tags users
//// @Accept json
//// @Produce json
//// @Param id path string true "user id"
//// @Success 200 {object} dto.HttpResponse[model.User]
//// @Failure 400 {object} dto.HttpResponse[any]
//// @Router /auth/info [get]
//func (handler *AuthHandler) Info(context *core.HttpContext) {
//	id := context.Gin.Param("id")
//	user, err := handler.service.FindById(id)
//	if err != nil {
//		context.Gin.Error(err)
//	} else {
//		context.JSON(http.StatusOK, dto.HttpResponse[*model.User]{
//			HttpStatus: http.StatusOK,
//			Data:       user,
//		})
//	}
//}

var AuthHandlerModule = fx.Options(fx.Provide(NewAuthHandler))
