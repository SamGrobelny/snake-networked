package main

import (
	"fmt"
	"net"
	"os"
)

type Player struct {
	addr *net.UDPAddr
	id   int
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	udpAdd, err := net.ResolveUDPAddr("udp", ":1337")

	handleErr(err)

	conn, err := net.ListenUDP("udp", udpAdd)

	fmt.Printf("listening on %s", conn.LocalAddr().String())

	handleErr(err)

	defer conn.Close()

	// var players []Player

	// players := make(map[*net.UDPAddr]Player, 0)

	for {
		var buffer [512]byte
		// for i, p := range players {
		//
		// }
		_, client_addr, err := conn.ReadFromUDP(buffer[0:])

		_, err = conn.WriteToUDP([]byte("pong\n"), client_addr)

		if err != nil {
			fmt.Println("error sending pong: ", err)
		}

		//TODO: handle the player set
		//TODO: choose whether to send update or send current game state

		handleErr(err)
	}
}
