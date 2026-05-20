package routes

import (
	"github.com/pascalallen/go-clean-arch/internal/app/application/http/action/auth"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
)

func (r Router) Auth(queryBus messaging.QueryBus, commandBus messaging.CommandBus) {
	v := r.engine.Group(v1)
	{
		a := v.Group("/auth")
		{
			a.POST("/register", auth.HandleRegisterUser(queryBus, commandBus))
			a.POST("/login", auth.HandleLoginUser(queryBus))
		}
	}
}
