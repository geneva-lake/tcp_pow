package service

import (
	"github.com/geneva-lake/tcp_pow/internal/protocol"
	"github.com/geneva-lake/tcp_pow/internal/vdf"
	"github.com/geneva-lake/tcp_pow/internal/wisdom"

	"github.com/google/uuid"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	wisdom     wisdom.Wisdomer //wisdom quote getter
	difficulty int32           // difficulty for vdf algorithm
}

func NewServer(wisdom wisdom.Wisdomer) *Server {
	s := Server{
		wisdom: wisdom,
	}
	return &s
}

func (s *Server) Init(cfg *Config) error {
	s.difficulty = cfg.Difficulty
	return s.wisdom.Init(cfg.WisdomFile)
}

// Server boot event processing
func (s *Server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng
	logging.Infof("server started")
	return
}

// Open connection processing
func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	conn := Connection{
		ID:        uuid.New(),
		wisdom:    s.wisdom,
		dificulty: s.difficulty,
		vdf:       vdf.NewVdf(),
	}
	c.SetContext(&conn)
	logging.Infof("connection opened id=%v addr=%s\n", conn.ID, c.RemoteAddr().String())
	return
}

// Close connection processing
func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	conn := c.Context().(*Connection)
	if err != nil {
		logging.Infof("connection closed id=%v addr=%s error=%v\n", conn.ID, c.RemoteAddr().String(), err)
	} else {
		logging.Infof("connection closed id=%v addr=%s\n", conn.ID, c.RemoteAddr().String())
	}
	return
}

// Upcoming message event processing
func (s *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	conn := c.Context().(*Connection)
	msgtype, data, err := protocol.NewCodec(c).Decode()
	if err != nil {
		logging.Infof("connection id=%v addr=%s error=%v\n", conn.ID, c.RemoteAddr().String(), err)
		return gnet.Close
	}
	msgtype, payload, err := conn.Process(msgtype, data)
	if err != nil {
		logging.Infof("connection id=%v addr=%s error=%v\n", conn.ID, c.RemoteAddr().String(), err)
		return gnet.Close
	}
	if msgtype == 0 {
		logging.Infof("connection id=%v addr=%s wrong message type", conn.ID, c.RemoteAddr().String())
		return gnet.Close
	}
	err = protocol.NewCodec(c).Encode(msgtype, payload)
	if err != nil {
		logging.Infof("connection id=%v addr=%s error=%v\n", conn.ID, c.RemoteAddr().String(), err)
		return gnet.Close
	}
	return gnet.None
}
