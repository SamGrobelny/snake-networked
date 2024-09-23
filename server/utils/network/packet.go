package network

import ()

type Packet struct {
	flag Flag
}

type JoinPacketClientServer struct {
	Packet    Packet
	SessionID uint64
}

type JoinPacketClientAck struct {
	Packet Packet
}

type LeavePacketClient struct {
	Packet    Packet
	SessionID uint64
}

type LeavePacketServerAck struct {
	Packet Packet
}

type PlayerStatePacketClient struct {
	Packet    Packet
	SessionID uint64
	Direction uint8
	Timestamp uint64
}

type PlayerPacket struct {
	Direction uint8
	NumPos    uint8
	Points    []Point
}

type GameStatePacketServer struct {
	Packet     Packet
	Timestamp  uint64
	GridWidth  uint8
	GridHeight uint8
	NumPlayers uint8
	Players    []PlayerPacket
}
