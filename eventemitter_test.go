package events

import (
	"runtime/debug"
	"testing"
)

func expect(t *testing.T, res bool, msgs ...interface{}) {
	if !res {
		debug.PrintStack()
		t.Error(msgs...)
	}
}

func TestSimpleEvent(t *testing.T) {
	msg := "simple message"
	ee := &EventEmitter{}
	ee.On("message", func(msg string) {
		t.Log("on:", msg)
	})
	t.Log("emitting:", msg)
	ee.Emit("message", msg)
}

func TestInterfaceFn(t *testing.T) {
	msg := "simple message"
	ee := &EventEmitter{}
	ee.On("message", func(msg interface{}) {
		t.Log("on:", msg)
	})
	t.Log("emitting:", msg)
	ee.Emit("message", msg)
}

func TestMultiHandler(t *testing.T) {
	msg := "simple message"
	ee := &EventEmitter{}
	ee.On("message", func(msg string) {
		t.Log("on (string):", msg)
	})
	ee.On("message", func(msg interface{}) {
		t.Log("on (interface):", msg)
	})
	t.Log("emitting:", msg)
	ee.Emit("message", msg)
}

type position struct {
	x, y float32
}

func TestHandlerByType(t *testing.T) {
	pos := new(position)
	pos.x = 4.5
	pos.y = 5.6

	ee := &EventEmitter{}
	ee.On("message", func(msg string) {
		t.Log("on (string):", msg)
	})
	ee.On("message", func(msg interface{}) {
		t.Log("on (interface):", msg)
	})
	ee.On("message", func(pos position) {
		t.Log("on (position):", pos.x, pos.y)
	})

	t.Log("emitting:", *pos)
	ee.Emit("message", *pos)
}

func TestHandlerByNArgs(t *testing.T) {
	pos := new(position)
	pos.x = 4.5
	pos.y = 5.6

	what := "car"
	ee := &EventEmitter{}
	ee.On("message", func(msg string) {
		t.Log("on (string):", msg)
	})
	ee.On("message", func(msg string, pos position) {
		t.Log("on (string, position):", msg, pos)
	})
	ee.On("message", func(msg string, pos position) string {
		t.Log("on (string, position) (string):", msg, pos)
		return msg
	})
	ee.On("message", func(msg string, pos position, ran interface{}) {
		t.Log("on (string, position interface{}):", msg, pos, ran)
	})

	t.Log("emitting:", what, *pos)
	ee.Emit("message", what, *pos)
}

func TestOnce(t *testing.T) {
	ee := &EventEmitter{}

	timesHelloEmitted := 0
	ee.Once("hello", func(a, b string) {
		t.Log("hello")
		timesHelloEmitted++
	})
	ee.Emit("hello", "a", "b")
	ee.Emit("hello", "a", "b")
	ee.Emit("hello", "a", "b")
	ee.Emit("hello", "a", "b")
	// expect(t, timesHelloEmitted == 1, "timesHelloEmitted =", timesHelloEmitted)

	remove := func() {
		expect(t, false, "once->foo should not be emitted!")
	}
	ee.Once("foo", remove)
	ee.RemoveListener("foo", remove)
	ee.Emit("foo")

	timesRecurseEmitted := 0
	ee.Once("e", func() {
		ee.Emit("e")
		t.Log("(1) timesRecurseEmitted++")
		timesRecurseEmitted++
	})
	ee.Once("e", func() {
		t.Log("(2) timesRecurseEmitted++")
		timesRecurseEmitted++
	})
	ee.Emit("e")
	// expect(t, timesRecurseEmitted == 2, "timesRecurseEmitted =", timesRecurseEmitted)
}

// Due to go concurrent feature, the ModifyInEmit test differs from
// the original midify-in-emit test in node.js events pacakge.
func TestModifiyInEmit(t *testing.T) {
	e := &EventEmitter{}

	var callback1, callback2, callback3 func()
	callback1 = func() {
		t.Log("callback1")
		e.On("foo", callback2)
		e.On("foo", callback3)
		e.RemoveListener("foo", callback1)
	}
	callback2 = func() {
		t.Log("callback2")
		e.RemoveListener("foo", callback2)
	}
	callback3 = func() {
		t.Log("callback3")
		e.RemoveListener("foo", callback3)
	}

	e.On("foo", callback1)
	e.Emit("foo")
	e.Emit("foo")
	e.Emit("foo")
	e.Emit("foo")

	e.On("foo", callback1)
	e.On("foo", callback2)
	e.RemoveAllListeners("foo")
	e.On("foo", callback2)
	e.On("foo", callback3)
	e.Emit("foo")
}
