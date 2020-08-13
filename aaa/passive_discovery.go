package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

const (
	udpBufferSize   = 2048
	pdListenAddress = "239.255.255.250:1982"
	pdTimeout 		= 30 * time.Second
)

func main() {
	fmt.Println("Warning: You may need to turn off your light bulbs via power switch and turn them back on " +
		"before discovery if you wont get the expected result.")
	fmt.Println()
	fmt.Println("Starting passive discovery...")

	var addr *net.UDPAddr
	var socket *net.UDPConn
	var err error

	addr, err = net.ResolveUDPAddr("udp", pdListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	socket, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully bound to multicast group. Waiting for periodic advertisements...")

	err = socket.SetReadBuffer(udpBufferSize)
	if err != nil {
		log.Fatal(err)
	}

	err = socket.SetReadDeadline(time.Now().Add(pdTimeout))
	if err != nil {
		log.Fatal(err)
	}

	for {
		var err error
		var buffer []byte

		buffer = make([]byte, udpBufferSize)
		_, _, err = socket.ReadFromUDP(buffer)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				fmt.Println("Timed out, exiting.")
				break
			} else {
				log.Fatal(err)
			}
		}

		locationRegex, _ := regexp.Compile(".*Location: yeelight://([0-9\\.:]+).*")
		locationMatch := locationRegex.FindStringSubmatch(string(buffer))
		if len(locationMatch) >= 1 {
			fmt.Println("Found light bulb listening at:", locationMatch[1])
		}
	}
}
