package model

type HAMessage interface {
	GetType() MsgType
}

type Message struct {
	Type MsgType `json:"type"`
}

func (msg *Message) GetType() MsgType {
	return msg.Type
}
