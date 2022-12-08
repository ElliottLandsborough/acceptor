package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	if _, err := c.Write([]byte(string("ok"))); err != nil {
		fmt.Printf("couldn't write into the conn: %+v\n", err)
	}

	if err := c.Close(); err != nil {
		fmt.Printf("couldn't close the conn: %+v\n", err)
		return
	}
}

func main() {
	fromPortFlag := flag.Int("from", 1024, "From port")
	tilPortFlag := flag.Int("until", 2048, "Until port")
	dieFlag := flag.Bool("die", true, "Die when port cannot be opened")
	flag.Parse()

	fromPort := *fromPortFlag
	tilPort := *tilPortFlag
	die := *dieFlag

	if fromPort < 1 {
		log.Panicf("From port must be between 1 and 65535")
	}

	if tilPort > 65535 {
		log.Panicf("From port must be between 1 and 65535")
	}

	log.Printf("spawning tcp servers in range from %d til %d", fromPort, tilPort)

	var toListen []net.Listener

	for i := 0; fromPort <= tilPort; i, fromPort = i+1, fromPort+1 {
		addr := ":" + strconv.FormatInt(int64(fromPort), 10)
		log.Printf("trying to bind %s", addr)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("can't create listener number {%d} - %+v", i, err)
			if die == true {
				goto toDie
			}

			continue
		}

		toListen = append(toListen, l)
	}

	for {
		for _, l := range toListen {
			c, err := l.Accept()
			if err != nil {
				fmt.Printf("can't accept the conn:%+v", err)
				return
			}
			go handleConnection(c)
		}
	}

toDie:
	for _, l := range toListen {
		if l != nil {
			if err := l.Close(); err != nil {
				fmt.Printf("can't close the listener:%+v", err)
				return
			}
		}
	}
}
