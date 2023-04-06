package service

import (
	"encoding/hex"
	"testing"

	"github.com/geneva-lake/tcp_pow/internal/protocol"

	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	client := NewClientEvents(&TestVdf{})
	input := []byte{0x00, 0x02, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x64, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	msgtype, payload, err := client.Process(protocol.Task, input)
	require.NoError(t, err, "error process task message")
	require.Equal(t, protocol.Solution, msgtype, "message type process task message")
	require.Equal(t, byte(0x80), payload[10], "payload process task message")
}

type TestVdf struct{}

func (v *TestVdf) Config(difficulty int, seed [32]byte) {}
func (v *TestVdf) Solve() [516]byte {
	inputVDF, _ := hex.DecodeString("0028f5de49d93dff7e2080a9bdadff1d63a2a4a143e6acedb814b78b49154ba6eb77d96d8c4ebefb2ae3f4b51af64219067c26693384eedffeca103767c2a4f4f0dd753a1e778aa372463f80a3fe01b2ca85a3be1707a8b82eeccffd0bc183a7f4c3c8854d3f46ec19bc797835e497b49db57b8a0fb0b87c3f3cfb3a631d12ee40ffe1bc410a72dd4804613e0bf6bf5968b75cbdc76ab45dae141b53645b9bfd5ffd667787b4941d1e1f306929844ced0fe90bf5e62632cb32e24f0f7dd276348dd3f769391da74456473513efd85b340f28504844b470187fdb5eccb9bf9e98897f1fba85f49f6fdbecaf6e18e12c34e4e525667f47de506cd5921ce818e026a06b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")
	var vdfBytes [516]byte
	copy(vdfBytes[:], inputVDF)
	return vdfBytes
}
func (v *TestVdf) Verify([516]byte) bool {
	return true
}
