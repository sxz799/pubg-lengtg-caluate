package main

import (
	"pubg-length-caluate/server"
	"pubg-length-caluate/web"
)

func main() {
	go server.Run()
	web.InitSocket()
}
