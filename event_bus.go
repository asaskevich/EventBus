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

// Subscribe - subscribe to a topic.
func (bus *EventBus) Subscribe(topic string, fn interface{}) {
	bus.lock.Lock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		bus.lock.Unlock()
		return
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = v
	bus.flagOnce[topic] = false
	bus.lock.Unlock()
}

// SubscribeOnce - subscribe to a topic once. Handler will be removed after executing.
func (bus *EventBus) SubscribeOnce(topic string, fn interface{}) {
	bus.lock.Lock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		bus.lock.Unlock()
		return
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = v
	bus.flagOnce[topic] = true
	bus.lock.Unlock()
}

// Unsubscribe - remove callback defined for a topic.
func (bus *EventBus) Unsubscribe(topic string) {
	bus.lock.Lock()
	if _, ok := bus.handlers[topic]; ok {
		delete(bus.handlers, topic)
	}
	bus.lock.Unlock()
}

func (bus *EventBus) PublishAsync(topic string, args ...interface{}) {
	bus.wg.Add(1)
	go func() {
		defer bus.wg.Done()
		bus.Publish(topic, args...)
	}()
}

// Publish - execute callback defined for a topic. Any addional argument will be tranfered to the callback.
func (bus *EventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock()
	if handler, ok := bus.handlers[topic]; ok {
		removeAfterExec, _ := bus.flagOnce[topic]
		args_ := make([]reflect.Value, 0)
		for _, arg := range args {
			args_ = append(args_, reflect.ValueOf(arg))
		}
		handler.Call(args_)
		if removeAfterExec {
			delete(bus.handlers, topic)
			bus.flagOnce[topic] = false
		}
	}
	bus.lock.Unlock()
}

func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
