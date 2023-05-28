package main

import (
	"fmt"
	"net"

	"github.com/julioguillermo/OXO/oxo"
	"github.com/julioguillermo/staticneurogenetic"
)

const (
	None = byte(iota)
	Reset
	Play
	Read
	Start
)

func main() {
	server, err := net.Listen("tcp", "0.0.0.0:9300")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err == nil {
			go process(conn)
		}
	}
}

func process(conn net.Conn) {
	ia, err := staticneurogenetic.LoadFromBin("oxo.bin")
	if err != nil {
		return
	}
	addr := conn.RemoteAddr().String()
	fmt.Printf("\033[32m%s:\033[0m open\n", addr)
	defer func() {
		fmt.Printf("\033[31m%s:\033[0m close\n", addr)
	}()
	game := oxo.NewOXO()
	ctl := make([]byte, 1)
	for {
		_, err := conn.Read(ctl)
		if err != nil {
			return
		}

		switch ctl[0] {
		case Reset:
			game.Reset()
		case Play:
			_, err = conn.Read(ctl)
			if err != nil {
				return
			}
			if ctl[0] < 9 {
				valid := game.Play(int(ctl[0]))
				if valid {
					_, pos := ia.MaxOutput(0, game.State())
					valid = game.Play(pos)
					for !valid {
						pos = (pos + 1) % 9
						valid = game.Play(pos)
					}
				}
			}
		case Read:
			_, err = conn.Write(game.Bytes())
			if err != nil {
				return
			}
		case Start:
			_, pos := ia.MaxOutput(0, game.State())
			valid := game.Play(pos)
			for !valid {
				pos = (pos + 1) % 9
				valid = game.Play(pos)
			}
		}
	}
}
