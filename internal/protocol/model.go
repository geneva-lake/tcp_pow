package protocol

type MessageType int16

// Protocol message types
const (
	Challenge          MessageType = 1
	Task               MessageType = 2
	Solution           MessageType = 3
	Citation           MessageType = 4
	VerificationFailed MessageType = 5
)
