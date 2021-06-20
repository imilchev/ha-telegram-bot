package model

type SubscribeEventsMessage struct {
	Message
	Id        int       `json:"id"`
	EventType EventType `json:"event_type"`
}

func NewSubscribeEventsMessage(id int, eventType EventType) SubscribeEventsMessage {
	return SubscribeEventsMessage{
		Message:   Message{Type: SubscribeEventsMsgType},
		Id:        id,
		EventType: eventType,
	}
}
