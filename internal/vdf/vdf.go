package vdf

import "github.com/harmony-one/vdf/src/vdf_go"

// Vdf algorithm interface
type VdfProcessor interface {
	Config(int, [32]byte)
	Solve() [516]byte
	Verify([516]byte) bool
}

type Vdf struct {
	difficulty int
	seed       [32]byte
}

func NewVdf() *Vdf {
	v := Vdf{}
	return &v
}

func (v *Vdf) Config(difficulty int, seed [32]byte) {
	v.difficulty = difficulty
	v.seed = seed
}

// Envelope vdf solution
func (v *Vdf) Solve() [516]byte {
	vdf := vdf_go.New(v.difficulty, v.seed)
	outputChannel := vdf.GetOutputChannel()
	vdf.Execute()
	output := <-outputChannel
	return output
}

// Envelope vdf verification
func (v *Vdf) Verify(solution [516]byte) bool {
	vdf := vdf_go.New(v.difficulty, v.seed)
	return vdf.Verify(solution)
}
