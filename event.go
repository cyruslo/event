package event

import (
	"sync"
)

type HandleFunc struct {
	method func(...interface{}) interface{}
	owner  int64
}

var eventMap sync.Map

func NewHandleFunc(owner int64, method func(...interface{}) interface{}) *HandleFunc {
	return &HandleFunc{
		method: method,
		owner:  owner,
	}
}

func Subscribe(event int, handle *HandleFunc) bool {
	if handle == nil {
		return false
	}
	if value, loaded := eventMap.Load(event); loaded {
		if handles, ok := value.(*sync.Map); ok {
			if _, ok := (*handles).Load(handle.owner); ok {
				return true
			} else {
				(*handles).Store(handle.owner, handle)
				eventMap.Store(event, handles)
			}
		} else {
			return false
		}
	} else {
		var handles = new(sync.Map)
		(*handles).Store(handle.owner, handle)
		eventMap.Store(event, handles)
	}
	return true
}

func UnSubscribe(event int, handle *HandleFunc) {
	if value, loaded := eventMap.Load(event); !loaded {
		return
	} else {
		if handles, ok := value.(*sync.Map); ok {
			if _, ok := (*handles).Load(handle.owner); ok {
				(*handles).Delete(handle.owner)
			}
			var count = 0
			(*handles).Range(func(k, v interface{}) bool {
				count++
				return true
			})
			if count <= 0 {
				eventMap.Delete(event)
			} else {
				eventMap.Store(event, handles)
			}
		} else {
			panic("unsubscribe event, loaded type err")
		}
	}
}

func Execute(event int, args ...interface{}) bool {
	if value, loaded := eventMap.Load(event); !loaded {
		return false
	} else {
		if handles, ok := value.(*sync.Map); ok {
			(*handles).Range(func(k, v interface{}) bool {
				if handle, r := v.(*HandleFunc); r {
					handle.method(args...)
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
