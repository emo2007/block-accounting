//go:build wireinject
// +build wireinject

package factory

import (
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/service"
	"github.com/google/wire"
)

func ProvideService(c config.Config) (service.Service, func(), error) {
	wire.Build(
		provideLogger,
		provideUsersRepository,
		provideUsersInteractor,
		provideJWTInteractor,
		interfaceSet,
		provideRestServer,
		service.NewService,
	)

	return &service.ServiceImpl{}, func() {}, nil
}
