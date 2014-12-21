package EventBus

import (
	"reflect"
	"sync"
)

// EventBus - box for handlers and callbacks.
type EventBus struct {
	handlers map[string]reflect.Value
	flagOnce map[string]bool
	lock     sync.Mutex
	wg       sync.WaitGroup
}

// New returns new EventBus with empty handlers.
func New() *EventBus {
	return &EventBus{
		make(map[string]reflect.Value),
		make(map[string]bool),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}

// Subscribe - subscribe to a channel.
func (bus *EventBus) Subscribe(channel string, fn interface{}) {
	bus.lock.Lock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		bus.lock.Unlock()
		return
	}
	v := reflect.ValueOf(fn)
	bus.handlers[channel] = v
	bus.flagOnce[channel] = false
	bus.lock.Unlock()
}

// SubscribeOnce - subscribe to a channel once. Handler will be removed after executing.
func (bus *EventBus) SubscribeOnce(channel string, fn interface{}) {
	bus.lock.Lock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		bus.lock.Unlock()
		return
	}
	v := reflect.ValueOf(fn)
	bus.handlers[channel] = v
	bus.flagOnce[channel] = true
	bus.lock.Unlock()
}

// Unsubscribe - remove callback defined for a channel.
func (bus *EventBus) Unsubscribe(channel string) {
	bus.lock.Lock()
	if _, ok := bus.handlers[channel]; ok {
		delete(bus.handlers, channel)
	}
	bus.lock.Unlock()
}

func (bus *EventBus) PublishAsync(channel string, args ...interface{}) {
	bus.wg.Add(1)
	go func() {
		defer bus.wg.Done()
		bus.Publish(channel, args...)
	}()
}

// Publish - execute callback defined for a channel. Any addional argument will be tranfered to the callback.
func (bus *EventBus) Publish(channel string, args ...interface{}) {
	bus.lock.Lock()
	if handler, ok := bus.handlers[channel]; ok {
		removeAfterExec, _ := bus.flagOnce[channel]
		args_ := make([]reflect.Value, 0)
		for _, arg := range args {
			args_ = append(args_, reflect.ValueOf(arg))
		}
		handler.Call(args_)
		if removeAfterExec {
			delete(bus.handlers, channel)
			bus.flagOnce[channel] = false
		}
	}
	bus.lock.Unlock()
}

func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
