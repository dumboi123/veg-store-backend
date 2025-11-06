package service

import (
	"fmt"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/domain/model"
	"veg-store-backend/internal/infrastructure/repository"

	"go.uber.org/fx"
)

type UserService interface {
	Name() string
	Start() error
	Stop() error

	Greeting() string
	FindById(id string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) Greeting() string {
	return "hello"
}

func (service *userService) FindById(id string) (*model.User, error) {
	if id == "1" {
		return nil, core.Error.NotFound.User

	} else {
		return &model.User{
			Name: "Test",
			Age:  18,
			Sex:  true,
		}, nil
	}
}

func (service *userService) FindByUsername(username string) (*model.User, error) {
	if username == "test" {
		return nil, core.Error.NotFound.User

	} else {
		return &model.User{
			Name: "Test",
			Age:  18,
			Sex:  true,
		}, nil
	}
}

func (service *userService) Name() string { return "UserService" }
func (service *userService) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", service.Name()))
	return nil
}
func (service *userService) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", service.Name()))
	return nil
}

var UserServiceModule = fx.Options(fx.Provide(NewUserService))

//func RegisterUserService(lifecycle fx.Lifecycle, service UserService) {
//	lifecycle.Append(fx.Hook{
//		OnStart: func(context context.Context) error {
//			return service.Start()
//		},
//		OnStop: func(context context.Context) error {
//			return service.Stop()
//		},
//	})
//}
//
//var UserServiceModule = fx.Options(
//	fx.Provide(NewUserService),
//	fx.Invoke(RegisterUserService),
//)
