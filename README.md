# 💡 🌈 Xiaomi Yeelight control tool written in Go

This is a simple tool which allow you to set colors on your Xiaomi Yeelight light bulbs in LAN network. **It requires you to enable LAN management in Yeelight settings.**

To run it:
----------
1. Install Golang - this depends on your operating system.
2. Execute `go build yeelight.go` to compile the source code. The result will be a binary file with the name `yeelight`.
3. Copy the executable to `/usr/local/bin` (on Unix systems) in order to be able to execute it from any location by running `sudo cp yeelight /usr/local/bin`.
4. Execute it with `yeelight <arguments here>`

Usage:
------
* `yeelight temperature <Light bulb IP address> <Light temperature in Kelvins 1700-6500> <Brightness 0-100>` - this command turns on the light bulb, switches it to white color mode and sets the desired color temperature and brightnes,
* `yeelight rgb <Light bulb IP address> <Red value 0-255> <Green value 0-255> <Blue value 0-255> <Brightness 0-100>` - this command turns on the light bulb, switches it to RGB/HSV color mode and sets the desired color and brightnes,
* `yeelight hsv <Light bulb IP address> <Hue 0-359> <Saturation 0-100> <Brightness/Value 0-100>` - this command turns on the light bulb, switches it to RGB/HSV color mode and sets the desired color and brightnes,
* `yeelight off <Light bulb IP address>` - this command turns off the light bulb,
* `yeelight help` - this command just displays help and exits.

Protocol details can be found here: https://www.yeelight.com/download/Yeelight_Inter-Operation_Spec.pdf
