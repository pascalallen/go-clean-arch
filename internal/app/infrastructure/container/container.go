package container

import (
	"database/sql"

	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/role"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/websocket"
)

type Container struct {
	DatabaseSession      *sql.DB
	Logger               logger.Logger
	PermissionRepository permission.Repository
	RoleRepository       role.Repository
	UserRepository       user.Repository
	CommandBus           messaging.CommandBus
	QueryBus             messaging.QueryBus
	EventDispatcher      messaging.EventDispatcher
	WebsocketHub         *websocket.Hub
}

func NewContainer(
	dbSession *sql.DB,
	logger logger.Logger,
	permissionRepo permission.Repository,
	roleRepo role.Repository,
	userRepo user.Repository,
	commandBus messaging.CommandBus,
	queryBus messaging.QueryBus,
	eventDispatcher messaging.EventDispatcher,
	websocketHub *websocket.Hub,
) Container {
	return Container{
		DatabaseSession:      dbSession,
		Logger:               logger,
		PermissionRepository: permissionRepo,
		RoleRepository:       roleRepo,
		UserRepository:       userRepo,
		CommandBus:           commandBus,
		QueryBus:             queryBus,
		EventDispatcher:      eventDispatcher,
		WebsocketHub:         websocketHub,
	}
}
