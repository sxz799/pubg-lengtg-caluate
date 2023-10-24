package main

import (
	"pubg-length-calculate/server"
	"pubg-length-calculate/web"
)

func main() {
	go server.Run()
	web.InitSocket()
}
