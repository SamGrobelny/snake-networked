package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"server/utils/game"
	"server/utils/network"
	"sync"
	"time"
)

// general packet format:
// | 1 byte: flags | packet specific format |
//
//	flags:
//		0x01 -	player state update
//		0x02 -	join
//		0x03 -	join ack
//		0x04 -	leave
//		0x05 -	leave ack
//
//
//	join packet:
//		client:
//			initial:	| flag: 0x02 |
//			ack:		| flag: 0x02 | 8 byte: session ID |
//		server:
//					| flag: 0x02 | 8 byte: session ID |
//
//	leave packet:
//		client:
//			initial:	| flag: 0x03 | 8 byte: session ID |
//		server:
//			ack:		| flag: 0x03 |
//
//	player state packet:
//		client:
//			send:		| flag: 0x01 | 8 byte: session ID | 1 byte: direction | timestamp |
//
//	initial game state send:
//		server:
//			send:		| flag: 0x01 | timestamp | 1 byte: grid width | 1 byte: grid height | 1 byte: num players | 1 byte * num players: direction | 1 byte * num players: num pos |
//					| num players * num pos * 2 bytes: x pos byte and y pos byte |
//		client:
//			ack:

const (
	gridWidth     = 32
	gridHeight    = 16
	maxPlayers    = 3
	afkTimeout    = 20 * time.Second
	tickDuration  = time.Second / 10
	serverAddress = ":1337"
)

const (
	None = iota
	Up
	Down
	Left
	Right
)

type Point struct {
	X uint8
	Y uint8
}

type Player struct {
	Address    *net.UDPAddr
	Id         uint32
	Direction  int
	Snake      []Point
	Lost       bool
	LastActive time.Time
}

type CustomPacket struct {
	PlayerId  uint32 // 4 bytes
	Timestamp int64  // 8 bytes
	Data      []byte // rest of stuff
}

// NOTE: I think this can be made better to try and prevent so many looks through the player arrays which takes a lot of time
// NOTE: Maybe make a matrix of the board that gets updated to just the matrix points around the player needs to be checked rather than every single player
// NOTE: If the players are scaled up, the would run very poorly (i think...)
type GameState struct {
	Players      map[uint32]*Player
	Mutex        sync.Mutex
	NextPlayerID uint32
	PlayerQueue  []*net.UDPAddr
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var gameState = GameState{Players: make(map[uint32]*Player)}

func resetGameState() {
	gameState.Players = make(map[uint32]*Player)
	gameState.PlayerQueue = []*net.UDPAddr{}
	gameState.NextPlayerID = 0
}

func encodePacket(playerId uint32, data []byte) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, playerId)
	binary.Write(&buffer, binary.BigEndian, time.Now().UnixMilli())
	buffer.Write(data)
	return buffer.Bytes()
}

func decodePacket(data []byte) CustomPacket {
	var packet CustomPacket
	buffer := bytes.NewReader(data)
	binary.Read(buffer, binary.BigEndian, &packet.PlayerId)
	binary.Read(buffer, binary.BigEndian, &packet.Timestamp)
	packet.Data = data[12:]
	return packet
}

func sendGameState(player *Player) {
	conn, err := net.DialUDP("udp", nil, player.Address)
	if err != nil {
		log.Println("error dialing udp:", err)
		return
	}

	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, uint32(len(gameState.Players)))
	for _, p := range gameState.Players {
		binary.Write(&buffer, binary.BigEndian, p.Id)
		binary.Write(&buffer, binary.BigEndian, len(p.Snake))

		for _, point := range p.Snake {
			binary.Write(&buffer, binary.BigEndian, uint32(point.X))
			binary.Write(&buffer, binary.BigEndian, uint32(point.Y))
		}
	}

	packet := encodePacket(player.Id, buffer.Bytes())
	conn.Write(packet)
}

func moveSnake(player *Player) {
	head := player.Snake[0]
	newHead := head

	switch player.Direction {
	case Up:
		newHead.Y = (newHead.Y - 1 + gridHeight) % gridHeight //wrap around if beyond bound
	case Down:
		newHead.Y = (newHead.Y + 1) % gridHeight
	case Left:
		newHead.X = (newHead.X - 1 + gridWidth) % gridWidth
	case Right:
		newHead.X = (newHead.X + 1) % gridWidth
	}

	player.Snake = append([]Point{newHead}, player.Snake...) //append new head in front
	player.Snake = player.Snake[:len(player.Snake)-1]        //remove last element
}

// NOTE: this can also be changed in the future to basically just move all the snakes and only check collision one per tick.
// NOTE: The logic here would be if a head ran into a body, whoevers head it is would lose.
// for now i just want the player to get deleted but in the future I want a add length to the person who "kills" them
func checkCollision(player *Player) {
	head := player.Snake[0]

	for _, enemy := range gameState.Players {
		for _, point := range enemy.Snake {
			if head == point {
				player.Lost = true
				return
			}
		}
	}
}

func gameLoop() {
	ticker := time.NewTicker(tickDuration)
	for range ticker.C {
		gameState.Mutex.Lock()

		if len(gameState.Players) == 0 {
			resetGameState()
			gameState.Mutex.Unlock()
			continue
		}

		for _, player := range gameState.Players {
			if player.Lost {
				continue //TODO: change this
			}
			moveSnake(player)
			checkCollision(player) //NOTE: move this below this loop in the future
		}

		for _, player := range gameState.Players {
			sendGameState(player)
		}

		gameState.Mutex.Unlock()
	}
}

func handlePacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	packet := decodePacket(data)
	gameState.Mutex.Lock()
	defer gameState.Mutex.Unlock()

	player, exists := gameState.Players[packet.PlayerId]

	if !exists {
		// if the max players are in the game, put in queue ( and do nothing else for now... )
		if len(gameState.Players) >= maxPlayers {
			gameState.PlayerQueue = append(gameState.PlayerQueue, addr)
			return
		}

		PlayerId := gameState.NextPlayerID
		gameState.NextPlayerID++

		var invalid map[Point]bool
		for _, p := range gameState.Players {
			for _, point := range p.Snake {
				invalid[point] = true
			}
		}

		//HACK: horrible way to do this. will change l8r
		position := Point{X: rand.IntN(gridWidth), Y: rand.IntN(gridHeight)}
		for !invalid[position] {
			position = Point{X: rand.IntN(gridWidth), Y: rand.IntN(gridHeight)}
		}

		player = &Player{
			Id:         PlayerId,
			Address:    addr,
			LastActive: time.Now(),
			Direction:  None,
			Snake:      []Point{position},
		}
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", serverAddress)
	handleErr(err)

	conn, err := net.ListenUDP("udp", addr)
	handleErr(err)

	defer conn.Close()

	go gameLoop()

	buffer := make([]byte, 1024)
	for {
		n, client_addr, err := conn.ReadFromUDP(buffer[0:])
		handleErr(err)

		handlePacket(conn, client_addr, buffer[:n])
	}
}
