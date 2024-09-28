package game

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
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
	LastActive time.Time
}

type GameState struct {
	Players      map[uint32]*Player
	DeadPlayers  map[uint32]*Player
	Mutex        sync.Mutex
	NextPlayerID string
}

func NewGameState() *GameState {
	gameState := &GameState{
		Players:      make(map[uint32]*Player),
		DeadPlayers:  make(map[uint32]*Player),
		NextPlayerID: generateID(),
	}

	return gameState
}

func (g *GameState) Reset() {
	g.Players = make(map[uint32]*Player)
	g.DeadPlayers = make(map[uint32]*Player)
	g.NextPlayerID = generateID()
}

func IncrementGameState(gameState *GameState) {
	for _, player := range gameState.Players {
		movePlayer(player)
	}
}

func generateID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func movePlayer(player *Player) {
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
// for now i just want the player to get deleted but in the future I want a add length to the person who "kills" them.
// Also change to move player to deadPlayers if "killed"
func checkCollision(player *Player, gameState *GameState) {
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
