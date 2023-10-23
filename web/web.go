package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pubg-length-caluate/server"
	"pubg-length-caluate/utils"
	"sync"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	server.ResultChannel <- "欢迎使用pubg火箭筒距离计算助手!"

	closeChannel := make(chan struct{})
	defer close(closeChannel)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading from WebSocket: %v", err)
				}
				<-closeChannel
				return
			}
		}
	}()

	for {
		select {
		case result := <-server.ResultChannel:
			err := conn.WriteMessage(websocket.TextMessage, []byte(result))
			if err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}
		case <-closeChannel:
			return
		}
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	server.ResultChannel <- "计算已开启!"
	server.MapOpen = true
}
func handleStop(w http.ResponseWriter, r *http.Request) {
	server.ResultChannel <- "计算已关闭!"
	server.MapOpen = false
}

func InitSocket() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("public"))) // 静态文件服务器
	http.HandleFunc("/start", handleStart)                // 静态文件服务器
	http.HandleFunc("/stop", handleStop)                  // 静态文件服务器
	ip, _ := utils.GetLocalIP()
	fmt.Println("服务器启动http://" + ip + ":3000/")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}

}
