package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	setPowerCmd       = "{\"id\": 1, \"method\": \"set_power\", \"params\":[\"%s\", \"smooth\", 500]}\r\n"
	setTemperatureCmd = "{\"id\": 1, \"method\": \"set_ct_abx\", \"params\":[%d, \"smooth\", 500]}\r\n"
	setRGBCmd         = "{\"id\": 1, \"method\": \"set_rgb\", \"params\":[%d, \"smooth\", 500]}\r\n"
	setHSVCmd         = "{\"id\": 1, \"method\": \"set_hsv\", \"params\":[%d, %d, \"smooth\", 500]}\r\n"
	setBrightnessCmd  = "{\"id\": 1, \"method\": \"set_bright\", \"params\":[%d, \"smooth\", 500]}\r\n"
	setTimeout        = 5 * time.Second
	udpBufferSize     = 2048
	adListenAddress   = ":0"
	adWriteAddress    = "239.255.255.250:1982"
	adMsg             = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\n" +
		"ST: wifi_bulb\r\n"
	adTimeout         = 10 * time.Second
	pdListenAddress   = "239.255.255.250:1982"
	pdTimeout         = 30 * time.Second
)

func main() {
	if len(os.Args) < 2 {
		displayHelp("You need to choice some action.")
	}

	switch command := os.Args[1]; command {
	case "temperature":
		temperature()
	case "rgb":
		rgb()
	case "hsv":
		hsv()
	case "off":
		off()
	case "active_discovery":
		activeDiscovery()
	case "passive_discovery":
		passiveDiscovery()
	case "help":
		displayHelp("")
	default:
		displayHelp("Bad action.")
	}
}

func temperature() {
	var host string
	var value int
	var brightness int
	var err error

	if len(os.Args) != 5 {
		displayHelp("Wrong number of parameters.")
	}

	if host = net.ParseIP(os.Args[2]).String(); host == "" {
		displayHelp("Wrong IP pdListenAddress format.")
	}

	if value, err = strconv.Atoi(os.Args[3]); err != nil || value < 1700 || value > 6500 {
		displayHelp("Wrong temperature value. It has to be between 1700 and 6500 Kelvins.")
	}

	if brightness, err = strconv.Atoi(os.Args[4]); err != nil || brightness < 0 || brightness > 100 {
		displayHelp("Wrong brightness value. It has to be between 0 and 100 percent.")
	}

	sendCommand(host, fmt.Sprintf(setPowerCmd, "on"))
	sendCommand(host, fmt.Sprintf(setTemperatureCmd, value))
	sendCommand(host, fmt.Sprintf(setBrightnessCmd, brightness))
}

func rgb() {
	var host string
	var red int
	var green int
	var blue int
	var brightness int
	var err error

	if len(os.Args) != 7 {
		displayHelp("Wrong number of parameters.")
	}

	if host = net.ParseIP(os.Args[2]).String(); host == "" {
		displayHelp("Wrong IP pdListenAddress format.")
	}

	if red, err = strconv.Atoi(os.Args[3]); err != nil || red < 0 || red > 255 {
		displayHelp("Wrong red value. It has to be between 0 and 255.")
	}

	if green, err = strconv.Atoi(os.Args[4]); err != nil || green < 0 || green > 255 {
		displayHelp("Wrong green value. It has to be between 0 and 255.")
	}

	if blue, err = strconv.Atoi(os.Args[5]); err != nil || blue < 0 || blue > 255 {
		displayHelp("Wrong blue value. It has to be between 0 and 255.")
	}

	if brightness, err = strconv.Atoi(os.Args[6]); err != nil || brightness < 0 || brightness > 100 {
		displayHelp("Wrong brightness value. It has to be between 0 and 100 percent.")
	}

	sendCommand(host, fmt.Sprintf(setPowerCmd, "on"))
	sendCommand(host, fmt.Sprintf(setRGBCmd, red<<16+green<<8+blue))
	sendCommand(host, fmt.Sprintf(setBrightnessCmd, brightness))
}

func hsv() {
	var host string
	var hue int
	var saturation int
	var brightness int
	var err error

	if len(os.Args) != 6 {
		displayHelp("Wrong number of parameters.")
	}

	if host = net.ParseIP(os.Args[2]).String(); host == "" {
		displayHelp("Wrong IP pdListenAddress format.")
	}

	if hue, err = strconv.Atoi(os.Args[3]); err != nil || hue < 0 || hue > 359 {
		displayHelp("Wrong hue value. It has to be between 0 and 359.")
	}

	if saturation, err = strconv.Atoi(os.Args[4]); err != nil || saturation < 0 || saturation > 255 {
		displayHelp("Wrong saturation value. It has to be between 0 and 255.")
	}

	if brightness, err = strconv.Atoi(os.Args[5]); err != nil || brightness < 0 || brightness > 100 {
		displayHelp("Wrong brightness value. It has to be between 0 and 100 percent.")
	}

	sendCommand(host, fmt.Sprintf(setPowerCmd, "on"))
	sendCommand(host, fmt.Sprintf(setHSVCmd, hue, saturation))
	sendCommand(host, fmt.Sprintf(setBrightnessCmd, brightness))
}

func off() {
	var host string

	if len(os.Args) != 3 {
		displayHelp("Wrong number of parameters.")
	}

	if host = net.ParseIP(os.Args[2]).String(); host == "" {
		displayHelp("Wrong IP pdListenAddress format.")
	}

	sendCommand(host, fmt.Sprintf(setPowerCmd, "off"))
}

func sendCommand(host string, command string) {
	var conn net.Conn
	var status string
	var err error

	conn, err = net.Dial("tcp", host + ":55443")
	if err != nil {
		fmt.Println("Error: Could not connect with the light bulb. Check IP pdListenAddress and make sure that " +
			"local connections are enabled in Yeelight settings.")
		os.Exit(1)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("Error: Can not close the connection with the light bulb.")
		}
	}()

	_, err = fmt.Fprintf(conn, command)
	if err != nil {
		fmt.Println("Error: Connection problem.")
		os.Exit(1)
	}

	err = conn.SetReadDeadline(time.Now().Add(setTimeout))
	if err != nil {
		fmt.Println("Error: Connection problem.")
		os.Exit(1)
	}

	status, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error: Connection problem.")
		os.Exit(1)
	}

	if !strings.Contains(status, "{\"id\":1,\"result\":[\"ok\"]}") {
		fmt.Println("Error: Received bad status code.\r\n" + status)
		os.Exit(1)
	}
}

func activeDiscovery() {
	fmt.Println("Xiaomi Yeelight WiFi light bulb CLI control tool written in Go")
	fmt.Println("==============================================================")
	fmt.Println()
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

func passiveDiscovery() {
	fmt.Println("Xiaomi Yeelight WiFi light bulb CLI control tool written in Go")
	fmt.Println("==============================================================")
	fmt.Println()
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
			fmt.Println("Got advertisement from light bulb listening at:", locationMatch[1])
		}
	}
}

func displayHelp(error string) {
	fmt.Println("Xiaomi Yeelight WiFi light bulb CLI control tool written in Go")
	fmt.Println("==============================================================")
	fmt.Println()

	if error != "" {
		fmt.Println("Error:\r\n\t" + error)
		fmt.Println()
	}

	fmt.Println("Usage:")
	fmt.Println("\tyeelight temperature <Light bulb IP pdListenAddress> <Light temperature in Kelvins 1700-6500> " +
		"<Brightness 0-100>")
	fmt.Println("\tyeelight rgb <Light bulb IP pdListenAddress> <Red value 0-255> <Green value 0-255> <Blue " +
		"value 0-255> <Brightness 0-100>")
	fmt.Println("\tyeelight hsv <Light bulb IP pdListenAddress> <Hue 0-359> <Saturation 0-100> <Brightness/Value " +
		"0-100>")
	fmt.Println("\tyeelight off <Light bulb IP pdListenAddress>")
	fmt.Println("\tyeelight active_discovery")
	fmt.Println("\tyeelight passive_discovery")
	fmt.Println("\tyeelight help")
	fmt.Println()
	fmt.Println("Author: Aleksander Kurczyk")
	fmt.Println("Sources: https://github.com/akurczyk/yeelight_cli")

	if error != "" {
		os.Exit(1)
	}
}
