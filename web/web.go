package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pubg-length-calculate/server"
	"pubg-length-calculate/utils"
	"strconv"
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
	server.CalculateOpen = true
}
func handleStop(w http.ResponseWriter, r *http.Request) {
	server.ResultChannel <- "计算已关闭!"
	server.CalculateOpen = false
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	//获取请求参数
	r.ParseForm()
	calBaseLength := r.Form.Get("calBaseLength")
	float, err := strconv.ParseFloat(calBaseLength, 64)
	if err != nil {
		w.Write([]byte("参数错误"))
	} else {
		server.CalBaseLength = float / 100 * 113
		w.Write([]byte("更新成功"))
	}

}

func InitSocket() {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir("public"))) // 静态文件服务器
	http.HandleFunc("/start", handleStart)                // 静态文件服务器
	http.HandleFunc("/stop", handleStop)                  // 静态文件服务器
	http.HandleFunc("/config", handleConfig)              // 静态文件服务器
	ip, _ := utils.GetLocalIP()
	fmt.Println("服务器启动,可使用手机打开http://" + ip + ":3000/进行使用!")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("服务器启动失败: ", err)
	}

}
