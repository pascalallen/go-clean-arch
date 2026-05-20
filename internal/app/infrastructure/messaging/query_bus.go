package messaging

import (
	"fmt"

	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

type Query interface {
	QueryName() string
}

type QueryHandler interface {
	Handle(query Query) (any, error)
}

type QueryBus interface {
	RegisterHandler(queryType string, handler QueryHandler)
	Fetch(qry Query) (any, error)
}

type SynchronousQueryBus struct {
	handlers map[string]QueryHandler
	logger   logger.Logger
}

func NewSynchronousQueryBus(logger logger.Logger) QueryBus {
	return &SynchronousQueryBus{
		handlers: make(map[string]QueryHandler),
		logger:   logger,
	}
}

func (bus *SynchronousQueryBus) RegisterHandler(queryType string, handler QueryHandler) {
	bus.logger.Info("registering query handler", "queryType", queryType)
	bus.handlers[queryType] = handler
}

func (bus *SynchronousQueryBus) Fetch(query Query) (any, error) {
	bus.logger.Info("fetching query", "queryName", query.QueryName())
	handler, found := bus.handlers[query.QueryName()]
	if !found {
		bus.logger.Warn("no handler registered for query type", "queryType", query.QueryName())
		return nil, fmt.Errorf("no handler registered for query type: %s", query.QueryName())
	}

	results, err := handler.Handle(query)
	if err != nil {
		bus.logger.Error("error calling query handler", "error", err, "queryType", query.QueryName())
		return nil, fmt.Errorf("error calling query handler: %s", err)
	}

	return results, nil
}
