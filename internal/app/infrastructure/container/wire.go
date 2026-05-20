//go:build wireinject
// +build wireinject

package container

import (
	"github.com/google/wire"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/database"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/logger/slog"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/repository"
)

func InitializeContainer() Container {
	wire.Build(
		NewContainer,
		slog.New,
		database.NewPostgresSession,
		repository.NewPostgresPermissionRepository,
		repository.NewPostgresRoleRepository,
		repository.NewPostgresUserRepository,
		messaging.NewChannelCommandBus,
		messaging.NewSynchronousQueryBus,
		messaging.NewChannelEventDispatcher,
	)
	return Container{}
}
