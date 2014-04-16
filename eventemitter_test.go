package events

import (
	"testing"
)

func TestSimpleEvent(t *testing.T) {
	msg := "simple message"
	ee := NewEventEmitter()
	ee.On("message", func(msg string) {
		t.Log("on:", msg)
	})
	t.Log("emitting:", msg)
	ee.Emit("message", msg)
}

func TestInterfaceFn(t *testing.T) {
	msg := "simple message"
	ee := NewEventEmitter()
	ee.On("message", func(msg interface{}) {
		t.Log("on:", msg)
	})
	t.Log("emitting:", msg)
	ee.Emit("message", msg)
}

func TestMultiHandler(t *testing.T) {
	msg := "simple message"
	ee := NewEventEmitter()
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

	ee := NewEventEmitter()
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
	ee := NewEventEmitter()
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
