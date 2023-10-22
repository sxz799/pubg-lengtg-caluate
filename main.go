package main

import (
	"LengthCal/server"
	"LengthCal/web"
)

func main() {
	go server.Run()
	web.InitSocket()
}
