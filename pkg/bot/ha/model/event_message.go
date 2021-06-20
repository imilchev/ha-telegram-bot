package model

type EventMessage struct {
	Message
	Id    int   `json:"id"`
	Event Event `json:"event"`
}

func NewEventMessage() *EventMessage {
	return &EventMessage{Message: Message{Type: EventMsgType}}
}

type Event struct {
	EventType EventType `json:"event_type"`
	Data      Data      `json:"data"`
}

type Data struct {
	EntityId string `json:"entity_id"`
	NewState State  `json:"new_state"`
	OldState State  `json:"old_state"`
}

type State struct {
	EntityId   string     `json:"entity_id"`
	State      string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	FriendlyName      string `json:"friendly_name"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
}
