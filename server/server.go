package server

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"math"
)

type Location struct {
	X float64
	Y float64
}

var ScreenDPI float64
var NeedCal bool

var ResultChannel chan string

func init() {
	ResultChannel = make(chan string, 10)
	x1, _ := robotgo.GetScreenSize()
	x2, _ := robotgo.GetScaleSize()
	ScreenDPI = float64(x2) / float64(x1)
	fmt.Println("程序启动成功,当前ScreenDPI:", ScreenDPI)
}

func Run() {
	var lastOperation Location
	mapOpen := false
	calLength := func(operation, lastOperation Location) {
		x := operation.X - lastOperation.X
		y := operation.Y - lastOperation.Y
		ResultChannel <- fmt.Sprintf("当前距离为 %f 米", math.Sqrt(x*x+y*y)/113*100)
	}
	robotgo.EventHook(hook.MouseDown, []string{}, func(event hook.Event) {
		var operation Location
		if event.Button == 2 && mapOpen {
			operation.X = float64(event.X) / ScreenDPI
			operation.Y = float64(event.Y) / ScreenDPI
			if NeedCal {
				calLength(operation, lastOperation)
			} else {
				ResultChannel <- "坐标1已获取,等待获取坐标2..."
			}
			lastOperation = operation
			NeedCal = !NeedCal
		}
	})
	robotgo.EventHook(hook.KeyDown, []string{}, func(event hook.Event) {
		if event.Keychar == 91 {
			mapOpen = true
			NeedCal = false
			ResultChannel <- "计算已开启!"
		}
		if event.Keychar == 93 {
			mapOpen = false
			ResultChannel <- "计算已关闭!"
		}
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}
