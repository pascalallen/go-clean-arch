package messaging

import (
	"sync"

	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

type Event interface {
	EventName() string
}

type Listener interface {
	Handle(event Event) error
}

type EventDispatcher interface {
	RegisterListener(eventType string, listener Listener)
	StartConsuming()
	Dispatch(evt Event)
	Shutdown()
}

type ChannelEventDispatcher struct {
	ch        chan Event
	listeners map[string]Listener
	logger    logger.Logger
	once      sync.Once
	wg        sync.WaitGroup
}

func NewChannelEventDispatcher(logger logger.Logger) EventDispatcher {
	return &ChannelEventDispatcher{
		ch:        make(chan Event, channelBufferSize),
		listeners: make(map[string]Listener),
		logger:    logger,
	}
}

func (e *ChannelEventDispatcher) RegisterListener(eventType string, listener Listener) {
	e.logger.Info("registering event listener", "eventType", eventType)
	e.listeners[eventType] = listener
}

func (e *ChannelEventDispatcher) Dispatch(evt Event) {
	e.logger.Info("dispatching event", "eventName", evt.EventName())
	e.ch <- evt
}

func (e *ChannelEventDispatcher) StartConsuming() {
	e.logger.Info("starting event dispatcher consumption")
	e.wg.Add(1)
	defer e.wg.Done()
	for evt := range e.ch {
		e.processEvent(evt)
	}
}

func (e *ChannelEventDispatcher) Shutdown() {
	e.once.Do(func() {
		close(e.ch)
	})
	e.wg.Wait()
}

func (e *ChannelEventDispatcher) processEvent(evt Event) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error("panic recovered during event processing", "panic", r, "eventType", evt.EventName())
		}
	}()

	e.logger.Info("processing event", "eventType", evt.EventName())

	listener, found := e.listeners[evt.EventName()]
	if !found {
		e.logger.Warn("no listener registered for event type", "eventType", evt.EventName())
		return
	}

	err := listener.Handle(evt)
	if err != nil {
		e.logger.Error("error calling listener for event", "error", err, "eventType", evt.EventName())
		return
	}

	e.logger.Info("event processed successfully", "eventType", evt.EventName())
}
