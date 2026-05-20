package command_handler

import (
	"fmt"

	"github.com/pascalallen/go-clean-arch/internal/app/application/command"
	"github.com/pascalallen/go-clean-arch/internal/app/application/event"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
)

type RegisterUserHandler struct {
	Logger          logger.Logger
	UserRepository  user.Repository
	EventDispatcher messaging.EventDispatcher
}

func (h RegisterUserHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.RegisterUser)
	if !ok {
		h.Logger.Error("invalid command type passed to RegisterUserHandler", "command", cmd)
		return fmt.Errorf("invalid command type passed to RegisterUserHandler: %v", cmd)
	}

	h.Logger.Info("handling RegisterUser command", "command_id", c.Id.String(), "email_address", c.EmailAddress)

	u := user.Register(c.Id, c.FirstName, c.LastName, c.EmailAddress)

	h.Logger.Info("persisting user to repository", "user_id", u.Id.String())
	if err := h.UserRepository.Add(u); err != nil {
		h.Logger.Error("user registration failed", "error", err, "user_id", u.Id.String())
		return fmt.Errorf("user registration failed: %s", err)
	}

	h.Logger.Info("dispatching UserRegistered event", "user_id", u.Id.String())
	h.EventDispatcher.Dispatch(&event.UserRegistered{
		Id:           c.Id,
		FirstName:    c.FirstName,
		LastName:     c.LastName,
		EmailAddress: c.EmailAddress,
	})

	h.Logger.Info("successfully handled RegisterUser command", "command_id", c.Id.String())

	return nil
}

type DeleteUserHandler struct {
	Logger         logger.Logger
	UserRepository user.Repository
}

func (h DeleteUserHandler) Handle(cmd messaging.Command) error {
	c, ok := cmd.(*command.DeleteUser)
	if !ok {
		h.Logger.Error("invalid command type passed to DeleteUserHandler", "command", cmd)
		return fmt.Errorf("invalid command type passed to DeleteUserHandler: %v", cmd)
	}

	h.Logger.Info("handling DeleteUser command", "user_id", c.UserId.String())

	u, err := h.UserRepository.GetById(c.UserId)
	if err != nil {
		h.Logger.Error("user not found", "error", err, "user_id", c.UserId.String())
		return fmt.Errorf("user not found: %s", c.UserId.String())
	}

	if u == nil {
		h.Logger.Error("user not found", "user_id", c.UserId.String())
		return fmt.Errorf("user not found: %s", c.UserId.String())
	}

	if err := h.UserRepository.Remove(u); err != nil {
		h.Logger.Error("failed to remove user", "error", err, "user_id", c.UserId.String())
		return fmt.Errorf("failed to remove user: %w", err)
	}

	h.Logger.Info("successfully handled DeleteUser command", "user_id", c.UserId.String())

	return nil
}
