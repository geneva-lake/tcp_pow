package protocol

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tg := NewTestGnet()
	codec := NewCodec(tg)
	data := []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x64, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	tg.Write(data)
	msgtype, payload, err := codec.Decode()
	require.NoError(t, err, "error decode message")
	require.Equal(t, Task, msgtype, "decode message type")
	require.Equal(t, payload[10], byte(0xde), "payload decode message")
}

func TestEncode(t *testing.T) {
	tg := NewTestGnet()
	codec := NewCodec(tg)
	data := []byte{0x00, 0x00, 0x00, 0x64, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	err := codec.Encode(Task, data)
	require.NoError(t, err, "error encode message")
	payload := make([]byte, 42)
	n, err := tg.Read(payload)
	require.NoError(t, err, "error reading message from a pipe")
	require.Equal(t, 42, n, "quantity of readed bytes from pipe")
	msgtype := binary.BigEndian.Uint16(payload[:typeSize])
	require.Equal(t, Task, MessageType(msgtype), "encoded message type")
	require.Equal(t, byte(0xde), payload[10], "encoded payload")
}

type TestGnet struct {
	buffer []byte
}

func NewTestGnet() *TestGnet {
	tg := &TestGnet{
		buffer: make([]byte, 100),
	}
	return tg
}

func (tg *TestGnet) Write(buf []byte) (int, error) {
	copy(tg.buffer, buf)
	return 0, nil
}

func (tg *TestGnet) Read(buf []byte) (int, error) {
	copy(buf, tg.buffer[:len(buf)])
	return len(buf), nil
}
