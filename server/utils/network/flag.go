package network

type Flag uint8

const (
	PlayerStateUpdate Flag = 0x01 + iota
	InitialGameState
	InitialGameStateAck
	Join
	JoinAck
	Leave
	LeaveAck
)
