package messaging

import (
	"sync"

	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(command Command) error
}

type CommandBus interface {
	RegisterHandler(commandType string, handler CommandHandler)
	StartConsuming()
	Execute(cmd Command) error
	Shutdown()
}

const channelBufferSize = 256

type ChannelCommandBus struct {
	ch       chan Command
	handlers map[string]CommandHandler
	logger   logger.Logger
	once     sync.Once
	wg       sync.WaitGroup
}

func NewChannelCommandBus(logger logger.Logger) CommandBus {
	return &ChannelCommandBus{
		ch:       make(chan Command, channelBufferSize),
		handlers: make(map[string]CommandHandler),
		logger:   logger,
	}
}

func (bus *ChannelCommandBus) RegisterHandler(commandType string, handler CommandHandler) {
	bus.logger.Info("registering command handler", "commandType", commandType)
	bus.handlers[commandType] = handler
}

func (bus *ChannelCommandBus) Execute(cmd Command) error {
	bus.logger.Info("executing command", "commandName", cmd.CommandName())
	bus.ch <- cmd
	return nil
}

func (bus *ChannelCommandBus) StartConsuming() {
	bus.logger.Info("starting command bus consumption")
	bus.wg.Add(1)
	defer bus.wg.Done()
	for cmd := range bus.ch {
		bus.processCommand(cmd)
	}
}

func (bus *ChannelCommandBus) Shutdown() {
	bus.once.Do(func() {
		close(bus.ch)
	})
	bus.wg.Wait()
}

func (bus *ChannelCommandBus) processCommand(cmd Command) {
	defer func() {
		if r := recover(); r != nil {
			bus.logger.Error("panic recovered during command processing", "panic", r, "commandType", cmd.CommandName())
		}
	}()

	bus.logger.Info("processing command", "commandType", cmd.CommandName())

	handler, found := bus.handlers[cmd.CommandName()]
	if !found {
		bus.logger.Warn("no handler registered for command type", "commandType", cmd.CommandName())
		return
	}

	err := handler.Handle(cmd)
	if err != nil {
		bus.logger.Error("error calling command handler", "error", err, "commandType", cmd.CommandName())
		return
	}

	bus.logger.Info("command processed successfully", "commandType", cmd.CommandName())
}
