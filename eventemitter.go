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
	once bool
}

type EventEmitterInt interface {
	On(event string, listener interface{})
	Once(event string, listener interface{})
	Emit(event string, args ...interface{})
	RemoveListener(event string, listener interface{})
	RemoveAllListeners(evs ...string)
}

type EventEmitter struct {
	lock      sync.RWMutex
	listeners listenerType
}

func getEventHandler(fn interface{}, once bool) (handler *eventHandler) {
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return nil
	}

	handler = new(eventHandler)
	handler.once = once
	handler.fn = fnValue
	handler.args = make([]reflect.Type, fnType.NumIn())

	for i := range handler.args {
		handler.args[i] = fnType.In(i)
	}

	return handler
}

func (ee *EventEmitter) addListener(event string, listener interface{}, once bool) {
	el := getEventHandler(listener, once)
	if el == nil {
		return
	}
	ee.lock.Lock()
	defer ee.lock.Unlock()
	if ee.listeners == nil {
		ee.listeners = make(listenerType)
	}
	ls, found := ee.listeners[event]
	if !found || ls == nil {
		ls = list.New()
		ee.listeners[event] = ls
	}
	ls.PushBack(el)
}

func (ee *EventEmitter) On(event string, listener interface{}) {
	ee.addListener(event, listener, false)
}

func (ee *EventEmitter) Once(event string, listener interface{}) {
	ee.addListener(event, listener, true)
}

func tryCall(el *eventHandler, args []interface{}) {
	defer func() {
		recover()
	}()
	if len(args) == len(el.args) {
		callArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			callArgs[i] = reflect.ValueOf(arg)
			if !reflect.TypeOf(arg).AssignableTo(el.args[i]) {
				return
			}
		}
		el.fn.Call(callArgs)
	}
}

func (ee *EventEmitter) Emit(event string, args ...interface{}) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	ls, found := ee.listeners[event]
	if found {
		for l := ls.Front(); l != nil; {
			next := l.Next()
			eh := l.Value.(*eventHandler)
			if eh.once {
				ls.Remove(l)
			}
			ee.lock.Unlock()
			tryCall(eh, args)
			ee.lock.Lock()
			l = next
		}
	}
}

// This implementation has a limit:
// It only uses reflect.ValueOf(listener).Pointer() to determine the uniqueness
// of a func, which is far from sufficient. This cannot handle func with a
// receiver obj.method, nor func that is a func closure with different binds.
func (ee *EventEmitter) RemoveListener(event string, listener interface{}) {
	var e *list.Element
	ee.lock.Lock()
	defer func() {
		// recover()
		ee.lock.Unlock()
	}()
	ptr := reflect.ValueOf(listener).Pointer()
	ls, found := ee.listeners[event]
	if found {
		for e = ls.Front(); e != nil; e = e.Next() {
			eh := e.Value.(*eventHandler)
			if eh.fn.Pointer() == ptr {
				ls.Remove(e)
				return
			}
		}
	}
}

func (ee *EventEmitter) RemoveAllListeners(evs ...string) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	if ee.listeners != nil {
		for _, ev := range evs {
			delete(ee.listeners, ev)
		}
	}
}
