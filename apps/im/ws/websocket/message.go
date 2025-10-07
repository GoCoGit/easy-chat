package websocket

type Message struct {
	Method string `json:"method"`
	FromId string `json:"fromId"`
	Data   any    `json:"data"`
}

func NewMessage(fromId string, data any) *Message {
	return &Message{
		FromId: fromId,
		Data:   data,
	}
}
