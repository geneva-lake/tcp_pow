package service

import (
	"encoding/hex"
	"testing"

	"github.com/geneva-lake/tcp_pow/internal/protocol"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var solution [516]byte
var quote []byte = []byte("Don’t seek happiness–create it. You don’t need life to go your way to be happy.")
var verified = true

func init() {
	inputVDF, _ := hex.DecodeString("0028f5de49d93dff7e2080a9bdadff1d63a2a4a143e6acedb814b78b49154ba6eb77d96d8c4ebefb2ae3f4b51af64219067c26693384eedffeca103767c2a4f4f0dd753a1e778aa372463f80a3fe01b2ca85a3be1707a8b82eeccffd0bc183a7f4c3c8854d3f46ec19bc797835e497b49db57b8a0fb0b87c3f3cfb3a631d12ee40ffe1bc410a72dd4804613e0bf6bf5968b75cbdc76ab45dae141b53645b9bfd5ffd667787b4941d1e1f306929844ced0fe90bf5e62632cb32e24f0f7dd276348dd3f769391da74456473513efd85b340f28504844b470187fdb5eccb9bf9e98897f1fba85f49f6fdbecaf6e18e12c34e4e525667f47de506cd5921ce818e026a06b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	copy(solution[:], inputVDF)
}

func TestProcess(t *testing.T) {
	conn := Connection{
		ID:        uuid.New(),
		wisdom:    &TestWisdom{},
		dificulty: 100,
		vdf:       &TestVdf{},
	}
	msgtype, payload, err := conn.Process(protocol.Challenge, nil)
	require.NoError(t, err, "error process challenge message")
	require.Equal(t, 36, len(payload), "length process challenge message")
	require.Equal(t, protocol.Task, msgtype, "message type process challenge message")

	msgtype, payload, err = conn.Process(protocol.Solution, solution[:])
	require.NoError(t, err, "error process solution message")
	require.Equal(t, len(quote), len(payload), "length process solution message")
	require.Equal(t, protocol.Citation, msgtype, "message type process solution message")

	verified = false
	msgtype, payload, err = conn.Process(protocol.Solution, solution[:])
	require.NoError(t, err, "error process wrong solution message")
	require.Equal(t, 0, len(payload), "length process wrong solution message")
	require.Equal(t, protocol.VerificationFailed, msgtype, "message type process wrong solution message")
}

type TestVdf struct{}

func (v *TestVdf) Config(difficulty int, seed [32]byte) {}
func (v *TestVdf) Solve() [516]byte {
	return solution
}
func (v *TestVdf) Verify([516]byte) bool {
	return verified
}

type TestWisdom struct{}

func (w *TestWisdom) Init(filename string) error {
	return nil
}
func (w *TestWisdom) GetQuote() []byte {
	return quote
}
