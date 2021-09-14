package event

import (
	"github.com/go-kratos/kratos/pkg/log"
	"runtime/debug"
	"sync"
)

type IEvent interface {
	GetEventType() EventType
}

type eventHandler struct {
}

type listener struct {
	eventType EventType
	owner     int64 //监听者唯一id
	method    func(event IEvent)
}

type eventQueueElem struct {
	event IEvent
	*listener
}

var (
	listenerMap sync.Map
	eventQueue  chan *eventQueueElem
)

func NewListener(eventType EventType, owner int64, method func(event IEvent)) *listener {
	if owner < 0 {
		return nil
	}
	return &listener{
		eventType: eventType,
		owner:     owner,
		method:    method,
	}
}

//Subscribe a listener
func Subscribe(listener *listener) bool {
	if listener == nil {
		return false
	}
	if value, loaded := listenerMap.Load(listener.eventType); loaded {
		if handles, ok := value.(*sync.Map); ok {
			if _, load := handles.LoadOrStore(listener.owner, listener); load {
				return true
			} else {
				listenerMap.Store(listener.eventType, handles)
			}
		} else {
			return false
		}
	} else {
		var handles = new(sync.Map)
		handles.Store(listener.owner, listener)
		listenerMap.Store(listener.eventType, handles)
	}
	return true
}

// owner = -1:delete all of eventType
func UnSubscribe(eventType EventType, owner int64) {
	if value, loaded := listenerMap.Load(eventType); !loaded {
		return
	} else {
		if owner == -1 {
			listenerMap.Delete(eventType)
			return
		}
		if handles, ok := value.(*sync.Map); ok {
			if _, ok := handles.Load(owner); ok {
				handles.Delete(owner)
			}
			var count = 0
			handles.Range(func(k, v interface{}) bool {
				count++
				return true
			})
			if count <= 0 {
				listenerMap.Delete(eventType)
			} else {
				listenerMap.Store(eventType, handles)
			}
		} else {
			panic("unsubscribe event, loaded type err")
		}
	}
}

func Execute(event IEvent) bool {
	eventType := event.GetEventType()
	if value, loaded := listenerMap.Load(eventType); !loaded {
		return false
	} else {
		if handles, ok := value.(*sync.Map); ok {
			(*handles).Range(func(k, v interface{}) bool {
				if listener, ok := v.(*listener); ok {
					eventQueue <- &eventQueueElem{
						event:    event,
						listener: listener,
					}
					return true
				} else {
					return false
				}
			})
		} else {
			panic("execute event, loaded type err")
		}
	}
	return true
}

func (self *eventHandler) run() {
	log.Info("event handler start")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Event handler Recover:%s", err)
				log.Error("Stack:%s", debug.Stack())
				self.run()
			}
		}()
		for {
			ele := <-eventQueue
			ele.method(ele.event)
		}
	}()
}

func init() {
	listenerMap = sync.Map{}
	eventQueue = make(chan *eventQueueElem, EVENT_CHANNEL_LEN)
	for i := 0; i < EXECUTOR_NUM; i++ {
		handler := &eventHandler{}
		handler.run()
	}
}

