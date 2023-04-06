package protocol

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()
	tg := &TestGnet{pipeReader, pipeWriter}
	codec := NewCodec(tg)
	data := []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x64, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	tg.Write(data)
	msgtype, payload, err := codec.Decode()
	require.NoError(t, err, "error decode message")
	require.Equal(t, Task, msgtype, "decode message type")
	require.Equal(t, payload[10], 0xde, "payload decode message")
}

func TestEncode(t *testing.T) {
	pipeReader, pipeWriter := io.Pipe()
	tg := &TestGnet{pipeReader, pipeWriter}
	codec := NewCodec(tg)
	data := []byte{0x00, 0x00, 0x00, 0x64, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	go tg.Write(data)
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
	*io.PipeReader
	*io.PipeWriter
}

func (tg *TestGnet) Write(buf []byte) (int, error) {
	go tg.PipeWriter.Write(buf)
	return 0, nil
}
