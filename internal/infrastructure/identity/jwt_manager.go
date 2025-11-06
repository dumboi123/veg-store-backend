package identity

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/infra_interface"
	"veg-store-backend/util"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type jwtManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTManager() (infra_interface.JWTManager, error) {
	// Set config path to .../.../keypair
	privateKeyPath := util.GetConfigPathFromGoMod("secrets/keypair") + core.Configs.JWT.PrivateKeyPath
	publicKeyPath := util.GetConfigPathFromGoMod("secrets/keypair") + core.Configs.JWT.PublicKeyPath
	core.Logger.Info(fmt.Sprintf("Private key path: %s", privateKeyPath))
	core.Logger.Info(fmt.Sprintf("Public key path: %s", publicKeyPath))

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		core.Logger.Fatal("error to read private key", zap.Error(err))
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		core.Logger.Fatal("error to read public key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		core.Logger.Fatal("error to parse private key", zap.Error(err))
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		core.Logger.Fatal("error to parse public key", zap.Error(err))
	}

	return &jwtManager{privateKey, publicKey}, nil
}

func (manager *jwtManager) Sign(isRefresh bool, userID string, roles ...string) (string, error) {
	var Expiration time.Duration
	var err error
	if isRefresh {
		Expiration, err = util.ParseDuration(core.Configs.JWT.RefreshDuration)
		if err != nil {
			core.Logger.Fatal("error to parse string to duration", zap.Error(err))
		}
	}
	claims := &infra_interface.JWTClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    core.Configs.JWT.ExpectedIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(manager.privateKey)
}

func (manager *jwtManager) Verify(tokenStr string) (*infra_interface.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &infra_interface.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return manager.publicKey, nil
	})
	if err != nil {
		core.Logger.Fatal("invalid token", zap.Error(err))
	}

	claims, ok := token.Claims.(*infra_interface.JWTClaims)
	if !ok || !token.Valid {
		core.Logger.Fatal("invalid claims", zap.Error(err))
	}
	return claims, nil
}

func (manager *jwtManager) Name() string { return "JWTManager" }
func (manager *jwtManager) Start() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", manager.Name()))
	return nil
}
func (manager *jwtManager) Stop() error {
	core.Logger.Debug(fmt.Sprintf("%s initialized", manager.Name()))
	return nil
}

var JWTManagerModule = fx.Options(fx.Provide(NewJWTManager))

//func RegisterJWTManager(lifecycle fx.Lifecycle, manager infra_interface.JWTManager) {
//	lifecycle.Append(fx.Hook{
//		OnStart: func(context context.Context) error {
//			return manager.Start()
//		},
//		OnStop: func(context context.Context) error {
//			return manager.Stop()
//		},
//	})
//}
//
//var JWTManagerModule = fx.Options(
//	fx.Provide(NewJWTManager),
//	fx.Invoke(RegisterJWTManager),
//)
