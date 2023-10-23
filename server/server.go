package server

import (
	"fmt"
	"math"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var MapOpen = false

var ResultChannel chan string

func Run() {
	ResultChannel = make(chan string, 10)
	type Point struct {
		X float64
		Y float64
	}
	x1, _ := robotgo.GetScreenSize()
	x2, _ := robotgo.GetScaleSize()
	screenDPI := float64(x2) / float64(x1)
	var NeedCal bool
	var lastPoint Point
	calLength := func(operation, lastOperation Point) {
		x := operation.X - lastOperation.X
		y := operation.Y - lastOperation.Y
		ResultChannel <- fmt.Sprintf("当前距离为 %f 米", math.Hypot(x, y)/113*100)
	}
	robotgo.EventHook(hook.MouseDown, []string{}, func(event hook.Event) {
		var point Point
		if event.Button == 2 && MapOpen {
			point.X = float64(event.X) / screenDPI
			point.Y = float64(event.Y) / screenDPI
			if NeedCal {
				calLength(point, lastPoint)
			} else {
				ResultChannel <- "坐标1已获取,等待获取坐标2..."
			}
			NeedCal = !NeedCal
			lastPoint = point
		}
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}
