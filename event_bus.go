package EventBus

import (
	"fmt"
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

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *EventBus) Subscribe(topic string, fn interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = v
	bus.flagOnce[topic] = false
	return nil
}

// SubscribeOnce subscribes to a topic once. Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *EventBus) SubscribeOnce(topic string, fn interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = v
	bus.flagOnce[topic] = true
	return nil
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (bus *EventBus) HasCallback(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic];
	return ok
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *EventBus) Unsubscribe(topic string) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if _, ok := bus.handlers[topic]; ok {
		delete(bus.handlers, topic)
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// PublishAsync executes callback defined for a topic asynchronously. Useful for slow callbacks.
// Any addional argument will be tranfered to the callback.
func (bus *EventBus) PublishAsync(topic string, args ...interface{}) {
	bus.wg.Add(1)
	go func() {
		defer bus.wg.Done()
		bus.Publish(topic, args...)
	}()
}

// Publish executes callback defined for a topic. Any addional argument will be tranfered to the callback.
func (bus *EventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if handler, ok := bus.handlers[topic]; ok {
		removeAfterExec, _ := bus.flagOnce[topic]
		passedArguments := make([]reflect.Value, 0)
		for _, arg := range args {
			passedArguments = append(passedArguments, reflect.ValueOf(arg))
		}
		handler.Call(passedArguments)
		if removeAfterExec {
			delete(bus.handlers, topic)
			bus.flagOnce[topic] = false
		}
	}
}

// WaitAsync waits for all async callbacks to complete
func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
