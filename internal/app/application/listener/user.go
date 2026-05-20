package listener

import (
	"fmt"

	"github.com/pascalallen/go-clean-arch/internal/app/application/event"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
)

type UserRegistration struct {
	Logger logger.Logger
}

func (l UserRegistration) Handle(evt messaging.Event) error {
	e, ok := evt.(*event.UserRegistered)
	if !ok {
		l.Logger.Error("invalid event type passed to UserRegistration listener", "event", evt)
		return fmt.Errorf("invalid event type passed to UserRegistration listener: %v", evt)
	}

	l.Logger.Info("user registered", "user_id", e.Id.String(), "email_address", e.EmailAddress)

	return nil
}
