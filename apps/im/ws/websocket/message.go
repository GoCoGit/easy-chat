package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string `json:"method"`
	FromId    string `json:"fromId"`
	Data      any    `json:"data"`
}

func NewMessage(fromId string, data any) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fromId,
		Data:      data,
	}
}
