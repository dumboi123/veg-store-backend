package service

import (
	"fmt"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/infra_interface"

	"go.uber.org/fx"
)

type AuthenticationService interface {
	Name() string
	Start() error
	Stop() error

	Tokens(request dto.SignInRequest) (*dto.Tokens, error)
}

type authenticationService struct {
	userService UserService
	jwtManager  infra_interface.JWTManager
}

func NewAuthenticationService(userService UserService, jwtManager infra_interface.JWTManager) AuthenticationService {
	return &authenticationService{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

func (service *authenticationService) Tokens(request dto.SignInRequest) (*dto.Tokens, error) {
	var err error
	user, err := service.userService.FindByUsername(request.Username)
	if err != nil {
		return nil, core.Error.Invalid.Username
	}

	accessToken, err := service.jwtManager.Sign(false, user.ID)
	if err != nil {
		return nil, core.Error.Auth.Unauthenticated
	}
	refreshToken, err := service.jwtManager.Sign(true, user.ID)
	if err != nil {
		return nil, core.Error.Auth.Unauthenticated
	}

	return &dto.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *authenticationService) Name() string { return "UserRepository" }
func (service *authenticationService) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", service.Name()))
	return nil
}
func (service *authenticationService) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", service.Name()))
	return nil
}

var AuthenticationServiceModule = fx.Options(fx.Provide(NewAuthenticationService))

//func RegisterAuthenticationService(lifecycle fx.Lifecycle, service AuthenticationService) {
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
//var AuthenticationServiceModule = fx.Options(
//	fx.Provide(NewAuthenticationService),
//	fx.Invoke(RegisterAuthenticationService),
//)
