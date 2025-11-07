package handler

import (
	"net/http"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/domain/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

// Hello godoc
// @Summary Anh trai say hi
// @Description Anh trai say gex
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.HttpResponse[string]
// @Router /user/hello [get]
func (handler *UserHandler) Hello(context *core.HttpContext) {
	context.JSON(http.StatusOK, gin.H{
		"message": context.T(handler.service.Greeting(), map[string]interface{}{
			"name":  "Ben",
			"Count": 1,
		}),
	})
}

// Details godoc
// @Summary User details
// @Description Get details of a user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} dto.HttpResponse[model.User]
// @Failure 400 {object} dto.HttpResponse[any]
// @Router /user/details/{id} [get]
func (handler *UserHandler) Details(context *core.HttpContext) {
	id := context.Gin.Param("id")
	user, err := handler.service.FindById(id)
	if err != nil {
		context.Gin.Error(err)
	} else {
		context.JSON(http.StatusOK, dto.HttpResponse[*model.User]{
			HttpStatus: http.StatusOK,
			Data:       user,
		})
	}
}

func (handler *UserHandler) HealthCheck(ctx *core.HttpContext) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (handler *UserHandler) GetAllUsers(ctx *core.HttpContext) {
	ctx.JSON(http.StatusOK, gin.H{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
}

var UserHandlerModule = fx.Options(fx.Provide(NewUserHandler))
