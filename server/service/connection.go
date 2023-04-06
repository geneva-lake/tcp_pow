package service

import (
	"crypto/rand"

	"github.com/geneva-lake/tcp_pow/internal/protocol"
	"github.com/geneva-lake/tcp_pow/internal/vdf"
	"github.com/geneva-lake/tcp_pow/internal/wisdom"

	"github.com/google/uuid"
)

// Information about certain connection
type Connection struct {
	ID        uuid.UUID
	dificulty int32            // difficulty for vdf algorithm
	seed      [32]byte         // seed for vdf algorithm
	wisdom    wisdom.Wisdomer  // wisdom quotes getter
	vdf       vdf.VdfProcessor // vdf algorithm interface
	verified  bool             // verified state of connection
}

// Processing messages from client
func (conn *Connection) Process(msgtype protocol.MessageType, data []byte) (protocol.MessageType, []byte, error) {
	switch msgtype {
	case protocol.Challenge:
		payload, err := generateTask(conn.dificulty)
		if err != nil {
			return 0, nil, err
		}
		conn.seed = *(*[32]byte)(payload[4:])
		return protocol.Task, payload, nil
	case protocol.Solution:
		conn.vdf.Config(int(conn.dificulty), conn.seed)
		verified := conn.vdf.Verify(*(*[516]byte)(data))
		if verified {
			conn.verified = true
			return protocol.Citation, conn.wisdom.GetQuote(), nil
		}
		return protocol.VerificationFailed, nil, nil
	default:
		return 0, nil, nil
	}
}

// Generate payload for task message
func generateTask(difficulty int32) ([]byte, error) {
	payload := make([]byte, 36)
	_, err := rand.Read(payload[4:])
	if err != nil {
		return nil, err
	}
	payload[3] = byte(difficulty >> 0)
	payload[2] = byte(difficulty >> 8)
	payload[1] = byte(difficulty >> 16)
	payload[0] = byte(difficulty >> 24)
	return payload, err
}
