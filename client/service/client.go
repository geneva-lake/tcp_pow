package service

import (
	"bytes"
	"encoding/binary"

	"github.com/geneva-lake/tcp_pow/internal/protocol"
	"github.com/geneva-lake/tcp_pow/internal/vdf"

	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
)

// Contains logic for message processing
type ClientEvents struct {
	*gnet.BuiltinEventEngine
	vdf  vdf.VdfProcessor
	stop chan interface{}
}

func NewClientEvents(v vdf.VdfProcessor) *ClientEvents {
	ev := ClientEvents{
		stop: make(chan interface{}),
		vdf:  v,
	}
	return &ev
}

func (ev ClientEvents) Stop() chan interface{} {
	return ev.stop
}

// Open connection processing
// sending challenge message
func (ev *ClientEvents) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	logging.Infof("connection opened addr=%s\n", c.RemoteAddr().String())
	err := protocol.NewCodec(c).Encode(protocol.Challenge, nil)
	if err != nil {
		logging.Infof("connection read addr=%s error=%v\n", c.RemoteAddr().String(), err)
		return nil, gnet.Close
	}
	return nil, gnet.None
}

// Close connection processing
func (ev *ClientEvents) OnClose(c gnet.Conn, err error) gnet.Action {
	defer func() {
		ev.stop <- 1
	}()
	if err != nil {
		logging.Infof("connection closed addr=%s error=%v\n", c.RemoteAddr().String(), err)
	} else {
		logging.Infof("connection closed addr=%s\n", c.RemoteAddr().String())
	}
	return gnet.None
}

// Upcoming message event processing
func (ev *ClientEvents) OnTraffic(c gnet.Conn) (action gnet.Action) {
	msgtype, data, err := protocol.NewCodec(c).Decode()
	if err != nil {
		logging.Infof("connection addr=%s error=%v\n", c.RemoteAddr().String(), err)
		return gnet.Close
	}
	msgtype, payload, err := ev.Process(msgtype, data)
	if err != nil {
		logging.Infof("connection addr=%s error=%v\n", c.RemoteAddr().String(), err)
		return gnet.Close
	}
	if msgtype == -1 {
		logging.Infof("connection addr=%s wrong message type", c.RemoteAddr().String())
		return gnet.Close
	}
	if msgtype != 0 {
		err := protocol.NewCodec(c).Encode(msgtype, payload)
		if err != nil {
			logging.Infof("connection addr=%s error=%v\n", c.RemoteAddr().String(), err)
			return gnet.Close
		}
	}
	return gnet.None
}

// Processing incoming messages
func (ev *ClientEvents) Process(msgtype protocol.MessageType, data []byte) (protocol.MessageType, []byte, error) {
	switch msgtype {
	case protocol.Task:
		solution, err := ev.solve(data)
		if err != nil {
			return 0, nil, err
		}
		return protocol.Solution, solution, nil
	case protocol.Citation:
		logging.Infof("wisdom from server: %s", data)
		return 0, nil, nil
	default:
		return -1, nil, nil
	}
}

// Solve pow task
func (ev *ClientEvents) solve(data []byte) ([]byte, error) {
	var difficulty int32
	buf := bytes.NewBuffer(data[:4])
	err := binary.Read(buf, binary.BigEndian, &difficulty)
	if err != nil {
		return nil, err
	}
	seed := *(*[32]byte)(data[4:])
	ev.vdf.Config(int(difficulty), seed)
	output := ev.vdf.Solve()
	return output[:], nil
}
