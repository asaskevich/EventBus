EventBus
======
Package EventBus is the little and lightweight eventbus for GoLang.

#### Installation
Make sure that Go is installed on your computer.
Type the following command in your terminal:

	go get github.com/asaskevich/EventBus

After it the package is ready to use.

#### Import package in your project
Add following line in your `*.go` file:
```go
import "github.com/asaskevich/EventBus"
```
If you unhappy to use long `EventBus`, you can do something like this:
```go
import (
	evbus "github.com/asaskevich/EventBus"
)
```

#### Example
```go
func calculator(a int, b int) {
	fmt.Printf("%d\n", a + b)
}

func main() {
	bus := EventBus.New();
	bus.Subscribe("main:calculator", calculator);
	bus.Publish("main:calculator", 20, 40);
	bus.Unsubscribe("main:calculator");
}
```

#### Implemented methods
* **New()**
* **Subscribe()**
* **SubscribeOnce()**
* **Unsubscribe()**
* **Publish()**
* **PublishAsync()**

#### New()
New returns new EventBus with empty handlers.
```go
bus := EventBus.New();
```

#### Subscribe(channel string, fn interface{})
Subscribe to a channel.
```go
func Handler() { ... }
...
bus.Subscribe("channel:handler", Handler)
```

#### SubscribeOnce(channel string, fn interface{})
Subscribe to a channel once. Handler will be removed after executing.
```go
func HelloWorld() { ... }
...
bus.SubscribeOnce("channel:handler", HelloWorld)
```

#### Unsubscribe(channel string)
Remove callback defined for a channel.
```go
bus.Unsubscribe("channel:handler");
```

#### Publish(channel string, args ...interface{})
Execute callback defined for a channel. Any addional argument will be tranfered to the callback.
```go
func Handler(str string) { ... }
...
bus.Subscribe("channel:handler", Handler)
...
bus.Publish("channel:handler", "Hello, World!");
```

#### PublishAsync(channel string, args ...interface{})
Execute callback defined for a channel asynchronously. Useful for slow callbacks.
```go
func slowCalculator(a, b int) {
	time.Sleep(3 * time.Second)
	fmt.Printf("%d\n", a + b)
}
...
bus := EventBus.New();
bus.Subscribe("main:slow_calculator", slowCalculator);

bus.Publish("main:slow_calculator", 20, 60); // synchronous execution means wait.
fmt.Println("I got blocked waiting")

bus.PublishAsync("main:slow_calculator", 30, 70);

fmt.Println("start: do some stuff while waiting for a result")
fmt.Println("end: do some stuff while waiting for a result") 

bus.WaitAsync(); // wait for all async callbacks to complete
bus.Unsubscribe("main:slow_calculator");
```

#### Support
If you do have a contribution for the package feel free to put up a Pull Request or open Issue.
