package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":1337")

	handleErr(err)

	conn, err := net.DialUDP("udp", nil, addr)

	handleErr(err)

	ticker := time.NewTicker(time.Second / 2)

	for {
		select {
		case t := <-ticker.C:
			_, err = conn.Write([]byte("ping"))
			handleErr(err)

			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			delta := time.Now().Sub(t).Milliseconds()
			fmt.Printf("%s : %d\n", string(data), delta)

		}
	}

}
