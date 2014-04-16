package events

import (
	"container/list"
	"reflect"
	"sync"
)

type listenerType map[string]*list.List
type eventHandler struct {
	fn   reflect.Value
	args []reflect.Type
}

type EventEmitter struct {
	sync.RWMutex
	listeners listenerType
}

func NewEventEmitter() (ee *EventEmitter) {
	ee = new(EventEmitter)
	ee.listeners = make(listenerType)
	return
}

func getEventHandler(fn interface{}) (handler *eventHandler) {
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return nil
	}

	handler = new(eventHandler)
	handler.fn = fnValue
	handler.args = make([]reflect.Type, fnType.NumIn())

	for i := range handler.args {
		handler.args[i] = fnType.In(i)
	}

	return handler
}

func (ee *EventEmitter) On(event string, listener interface{}) {
	el := getEventHandler(listener)
	if el == nil {
		return
	}
	ee.Lock()
	defer ee.Unlock()
	ls, found := ee.listeners[event]
	if !found {
		ls = list.New()
		ee.listeners[event] = ls
	}
	ls.PushBack(el)
}

func tryCall(el *eventHandler, args []interface{}) {
	if len(args) == len(el.args) {
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
			if !reflect.TypeOf(arg).AssignableTo(el.args[i]) {
				return
			}
		}
		go el.fn.Call(callArgs)
	}
}

func (ee *EventEmitter) Emit(event string, args ...interface{}) {
	ee.RLock()
	defer ee.RUnlock()
	ls, found := ee.listeners[event]
	if found {
		for l := ls.Front(); l != nil; l = l.Next() {
			tryCall(l.Value.(*eventHandler), args)
		}
	}
}
