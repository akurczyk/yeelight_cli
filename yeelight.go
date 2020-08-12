package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	setPowerCmd = "{\"id\": 1, \"method\": \"set_power\", \"params\":[\"%s\", \"smooth\", 500]}\r\n"
	setTemperatureCmd = "{\"id\": 1, \"method\": \"set_ct_abx\", \"params\":[%d, \"smooth\", 500]}\r\n"
	setRGBCmd = "{\"id\": 1, \"method\": \"set_rgb\", \"params\":[%d, \"smooth\", 500]}\r\n"
	setHSVCmd = "{\"id\": 1, \"method\": \"set_hsv\", \"params\":[%d, %d, \"smooth\", 500]}\r\n"
	setBrightnessCmd = "{\"id\": 1, \"method\": \"set_bright\", \"params\":[%d, \"smooth\", 500]}\r\n"
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
		displayHelp("Wrong IP address format.")
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
		displayHelp("Wrong IP address format.")
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
	sendCommand(host, fmt.Sprintf(setRGBCmd, red<<16 + green<<8 + blue))
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
		displayHelp("Wrong IP address format.")
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
		displayHelp("Wrong IP address format.")
	}

	sendCommand(host, fmt.Sprintf(setPowerCmd, "off"))
}

func sendCommand(host string, command string) {
	var conn net.Conn
	var status string
	var err error

	conn, err = net.Dial("tcp", host+":55443")
	if err != nil {
		fmt.Println("Error: Could not connect with the light bulb. Check IP address and make sure that local connections are enabled in Yeelight settings.")
		os.Exit(1)
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, command)
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

func displayHelp(error string) {
	fmt.Println("Xiaomi Yeelight WiFi light bulb CLI control tool written in Go")
	fmt.Println("==============================================================")
	fmt.Println()

	if error != "" {
		fmt.Println("Error:\r\n\t" + error)
		fmt.Println()
	}

	fmt.Println("Usage:")
	fmt.Println("\tyeelight temperature <Light bulb IP address> <Light temperature in Kelvins 1700-6500> <Brightness 0-100>")
	fmt.Println("\tyeelight rgb <Light bulb IP address> <Red value 0-255> <Green value 0-255> <Blue value 0-255> <Brightness 0-100>")
	fmt.Println("\tyeelight hsv <Light bulb IP address> <Hue 0-359> <Saturation 0-100> <Brightness/Value 0-100>")
	fmt.Println("\tyeelight off <Light bulb IP address>")
	fmt.Println("\tyeelight help")
	fmt.Println()
	fmt.Println("Author: Aleksander Kurczyk")
	fmt.Println("Sources: https://github.com/akurczyk/yeelight_cli")
	fmt.Println("License: Creative Commons")

	if error != "" {
		os.Exit(1)
	}
}
