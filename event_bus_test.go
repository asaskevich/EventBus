package EventBus

import "testing"

func TestNew(t *testing.T) {
	bus := New();
	if bus == nil {
		t.Log("New EventBus not created!");
		t.Fail();
	}
}

func TestHasCallback(t *testing.T) {
	bus := New();
	bus.Subscribe("topic", func() {} );
	if bus.HasCallback("topic_topic") {
		t.Fail();
	}
	if !bus.HasCallback("topic") {
		t.Fail();
	}
}

func TestSubscribe(t *testing.T) {
	bus := New();
	if bus.Subscribe("topic", func() {} ) != nil {
		t.Fail();
	}
	if bus.Subscribe("topic", "String" ) == nil {
		t.Fail();
	}
}

func TestSubscribeOnce(t *testing.T) {
	bus := New();
	if bus.SubscribeOnce("topic", func() {} ) != nil {
		t.Fail();
	}
	if bus.SubscribeOnce("topic", "String" ) == nil {
		t.Fail();
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New();
	bus.Subscribe("topic", func() {} );
	if bus.Unsubscribe("topic") != nil {
		t.Fail();
	}
	if bus.Unsubscribe("topic") == nil {
		t.Fail();
	}
}

func TestPublish(t *testing.T) {
	bus := New();
	bus.Subscribe("topic", func(a int, b int) {
			if (a != b) {
				t.Fail();
			}
		} );
	bus.Publish("topic", 10, 10);
}
