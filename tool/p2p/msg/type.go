package msg

type MsgType uint16

const (
	MsgLogin MsgType = iota + 1
	MsgLogout
	MsgList
	MsgPunch
	MsgPing
	MsgPong
	MsgReply
	MsgText
	MsgEnd
)
