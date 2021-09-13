package event

import (
	"testing"
)

func Test(t *testing.T) {
	testEventMode()
}

func handleFunc(event IEvent) {
}

func testEventMode() {
	initId := 0
	for i := 0; i < 900000; i++ {
		listener := NewListener(EVENT_USER_LOGIN, int64(initId+i), handleFunc)
		Subscribe(listener)
	}
	//time.Sleep(time.Second * 20)
	event := &UserLoginEvent{}
	Execute(event)
}
