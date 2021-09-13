package event

type EventType int32

const EVENT_CHANNEL_LEN = 128

const EXECUTOR_NUM = 16

const (
	//Event
	EVENT_USER_LOGIN EventType = iota + 1
)

type UserLoginEvent struct {

}

func (self *UserLoginEvent) GetEventType() EventType {
	return EVENT_USER_LOGIN
}
