package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	typeSize = 2
	bodySize = 4
)

type Codec struct {
	c io.ReadWriter
}

func NewCodec(c io.ReadWriter) *Codec {
	return &Codec{c}
}

// Protocol format:
//
// * 0           2                       6
// * +-----------+-----------------------+
// * |   type    |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+

// Encode message
func (codec *Codec) Encode(t MessageType, buf []byte) error {
	bodyOffset := typeSize + bodySize
	msgLen := bodyOffset
	if buf != nil {
		msgLen = bodyOffset + len(buf)
	}
	data := make([]byte, msgLen)

	binary.BigEndian.PutUint16(data, uint16(t))
	if buf != nil {
		binary.BigEndian.PutUint32(data[typeSize:bodyOffset], uint32(len(buf)))
		copy(data[bodyOffset:msgLen], buf)
	}
	_, err := codec.c.Write(data)
	return err
}

// Decode message
func (codec *Codec) Decode() (MessageType, []byte, error) {
	bodyOffset := typeSize + bodySize
	buf := make([]byte, bodyOffset)
	n, err := codec.c.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	if n < bodyOffset {
		return 0, nil, errors.New("read error")
	}
	msgtype := binary.BigEndian.Uint16(buf[:typeSize])
	bodyLen := binary.BigEndian.Uint32(buf[typeSize:bodyOffset])
	if bodyLen == 0 {
		return MessageType(msgtype), nil, nil
	}
	buf = make([]byte, bodyLen)
	n, err = codec.c.Read(buf)
	if err != nil {
		return 0, nil, err
	}
	if n < int(bodyLen) {
		return 0, nil, errors.New("read error")
	}
	return MessageType(msgtype), buf, nil
}
