package rtmp

type MessageType byte

const (
	MsgSetChunkSize              MessageType = 1
	MsgAbortMessage              MessageType = 2
	MsgAcknowledgement           MessageType = 3
	MsgUserControl               MessageType = 4 //
	MsgWindowAcknowledgementSize MessageType = 5
	MsgSetPeerBandwidth          MessageType = 6
	MsgCommandAmf0               MessageType = 20
	MsgCommandAmf3               MessageType = 17
	MsgDataAmf0                  MessageType = 18
	MsgDataAmf3                  MessageType = 15
	MsgSharedObjectAmf0          MessageType = 19
	MsgSharedObjectAmf3          MessageType = 16
	MsgAudio                     MessageType = 8
	MsgVideo                     MessageType = 9
	MsgAggregate                 MessageType = 22
)

type MessageHeader struct {
}

type Message interface {
}

type ChunkSize struct {
}

type AbortMessage struct {
}

type Acknowledgement struct {
}

// type EventType uint16

// const (
// 	EventStreamBegin EventType = iota
// 	EventStreamEOF
// 	EventStreamDry
// 	EventSetBufferLength
// 	EventStreamIsRecorded
// 	_
// 	EventPingRequest
// 	EventPingResponse
// )

// type UserControl struct {
// }

type WindowAcknowledgementSize struct {
}

type LimitType uint8

const (
	LimitHard LimitType = iota
	LimitSoft
	LimitDynamic
)

type SetPeerBandwidth struct {
}

type Command interface {
	Message
}

type Data struct {
}

type SharedObject struct {
}

type Audio struct {
}

type Video struct {
}

type Aggregate struct {
}
