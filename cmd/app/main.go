package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pascalallen/go-clean-arch/internal/app/application/command"
	"github.com/pascalallen/go-clean-arch/internal/app/application/command_handler"
	"github.com/pascalallen/go-clean-arch/internal/app/application/event"
	"github.com/pascalallen/go-clean-arch/internal/app/application/listener"
	"github.com/pascalallen/go-clean-arch/internal/app/application/query"
	"github.com/pascalallen/go-clean-arch/internal/app/application/query_handler"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/container"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/database"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/routes"
)

func main() {
	serviceContainer := container.InitializeContainer()

	database.RunMigrations(serviceContainer.DatabaseSession)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		runConsumers(serviceContainer)
		configureServer(serviceContainer)
	}()

	<-stop
	fmt.Println("\nShutting down gracefully...")
	serviceContainer.CommandBus.Shutdown()
	serviceContainer.EventDispatcher.Shutdown()
}

func runConsumers(c container.Container) {
	setupCommandHandlers(c.CommandBus, c)
	setupEventListeners(c.EventDispatcher, c)
	setupQueryHandlers(c.QueryBus, c)

	go c.CommandBus.StartConsuming()
	go c.EventDispatcher.StartConsuming()
}

func setupCommandHandlers(commandBus messaging.CommandBus, c container.Container) {
	commandBus.RegisterHandler(command.RegisterUser{}.CommandName(), command_handler.RegisterUserHandler{
		Logger:          c.Logger,
		UserRepository:  c.UserRepository,
		EventDispatcher: c.EventDispatcher,
	})
	commandBus.RegisterHandler(command.DeleteUser{}.CommandName(), command_handler.DeleteUserHandler{
		Logger:         c.Logger,
		UserRepository: c.UserRepository,
	})
}

func setupEventListeners(eventDispatcher messaging.EventDispatcher, c container.Container) {
	eventDispatcher.RegisterListener(event.UserRegistered{}.EventName(), listener.UserRegistration{
		Logger: c.Logger,
	})
}

func setupQueryHandlers(queryBus messaging.QueryBus, c container.Container) {
	queryBus.RegisterHandler(query.ListUsers{}.QueryName(), query_handler.ListUsersHandler{
		Logger:         c.Logger,
		UserRepository: c.UserRepository,
	})
	queryBus.RegisterHandler(query.GetUserById{}.QueryName(), query_handler.GetUserByIdHandler{
		Logger:         c.Logger,
		UserRepository: c.UserRepository,
	})
	queryBus.RegisterHandler(query.GetUserByEmailAddress{}.QueryName(), query_handler.GetUserByEmailAddressHandler{
		Logger:         c.Logger,
		UserRepository: c.UserRepository,
	})
}

func configureServer(c container.Container) {
	gin.SetMode(os.Getenv("GIN_MODE"))

	router := routes.NewRouter()
	router.UseLogger(c.Logger)
	router.Auth(c.QueryBus, c.CommandBus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Serve(":" + port)
}
