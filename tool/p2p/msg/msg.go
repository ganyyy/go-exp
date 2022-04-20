package msg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"
)

const (
	HeadSize = int(unsafe.Sizeof(MsgHead{}))

	MagicNum = uint16(0x8976)
)

var (
	ErrMagicNum   = errors.New("the magic num not match")
	ErrHeadLength = errors.New("the head length not valid")
	ErrDataLength = errors.New("the data length not match")
)

type MsgHead struct {
	Magic  uint16
	Type   MsgType
	Length uint16
}

type Msg struct {
	MsgHead
	Data []byte
}

func WriteUint16(val uint16) []byte {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], val)
	return buf[:]
}

func ReadUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

func (m Msg) Pack() []byte {
	var buf = bytes.NewBuffer(make([]byte, 0, HeadSize+len(m.Data)))
	m.Magic = MagicNum
	buf.Write(WriteUint16(m.Magic)[:])
	buf.Write(WriteUint16(uint16(m.Type))[:])
	m.Length = uint16(len(m.Data))
	buf.Write(WriteUint16(m.Length)[:])
	buf.Write(m.Data)
	return buf.Bytes()
}

func (m *Msg) Unpack(data []byte) error {
	if len(data) < HeadSize {
		return ErrHeadLength
	}
	var reader = bytes.NewReader(data)

	var read = func() uint16 {
		var buf [2]byte
		reader.Read(buf[:])
		return ReadUint16(buf[:])
	}

	m.Magic = read()
	if m.Magic != MagicNum {
		return ErrMagicNum
	}
	m.Type = MsgType(read())
	m.Length = read()
	if m.Length != uint16(reader.Len()) {
		return ErrDataLength
	}
	m.Data = make([]byte, m.Length)
	reader.Read(m.Data)
	return nil
}

func NewMsgNoData(msgType MsgType) Msg {
	return Msg{
		MsgHead: MsgHead{
			Magic: MagicNum,
			Type:  msgType,
		},
	}
}

func NewMsg[T ~string | []byte](msgType MsgType, data T) Msg {
	return Msg{
		MsgHead: MsgHead{
			Magic:  MagicNum,
			Type:   msgType,
			Length: uint16(len(data)),
		},
		Data: []byte(data),
	}
}

func ReadMsg(buf []byte) (msg Msg, err error) {
	err = msg.Unpack(buf)
	return
}
