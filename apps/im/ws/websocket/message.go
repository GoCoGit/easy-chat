package websocket

import "time"

type FrameType uint8

const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9
)

type Message struct {
	FrameType `json:"frameType"`
	Id        string    `json:"id"`
	AckSeq    int       `json:"ackSeq"`
	AckTime   time.Time `json:"ackTime"`
	ErrCount  int       `json:"errCount"`
	Method    string    `json:"method"`
	FromId    string    `json:"fromId"`
	Data      any       `json:"data"`
}

func NewMessage(fromId string, data any) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fromId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
