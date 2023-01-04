package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ribincao/ribin-game-server/timer"
)

func TestTimerWheel(t *testing.T) {
	tw, _ := timer.NewTimeWheel(time.Millisecond*100, 4, timer.TickSafeMode())
	tw.Start()
	defer tw.Stop()

	i := 0
	tw.AddTask("Circle", time.Second*1, func() {
		fmt.Println("-", time.Now().Second(), i)
		i++
	}, timer.CIRCLE_MODE)

	tw.AddTask("NotCircle", time.Second*10, func() {
		fmt.Println("+", time.Now().Second(), "end")
	}, timer.NOT_CIRCLE_MODE)

	time.Sleep(time.Second * 15)
}
