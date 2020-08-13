package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

const (
	udpBufferSize   = 2048
	adListenAddress = ":0"
	adWriteAddress  = "239.255.255.250:1982"
	adMsg           = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb\r\n"
	adTimeout       = 10 * time.Second
)

func main() {
	fmt.Println("Warning: You may need to turn off your light bulbs via power switch and turn them back on " +
		"before discovery if you wont get the expected result.")
	fmt.Println()
	fmt.Println("Starting active discovery...")

	var listenAddr, writeAddr *net.UDPAddr
	var socket *net.UDPConn
	var err error

	listenAddr, err = net.ResolveUDPAddr("udp", adListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	writeAddr, err = net.ResolveUDPAddr("udp", adWriteAddress)
	if err != nil {
		log.Fatal(err)
	}

	socket, err = net.ListenUDP("udp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	err = socket.SetReadBuffer(udpBufferSize)
	if err != nil {
		log.Fatal(err)
	}

	_, err = socket.WriteToUDP([]byte(adMsg), writeAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Multicast packet sent. Waiting for responses...")

	for {
		var err error
		var buffer []byte

		err = socket.SetReadDeadline(time.Now().Add(adTimeout))
		if err != nil {
			log.Fatal(err)
		}

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
