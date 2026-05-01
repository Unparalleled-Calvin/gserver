package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer, err := NewTimer(time.Second, time.Hour)
	if err != nil {
		t.Errorf("创建定时器失败: %v", err)
	}
	timer.AddTask(&Task{
		func() { fmt.Println("hello world") },
		time.Second * 59,
	})
	timer.AddTask(&Task{
		func() { fmt.Println("hello world") },
		time.Second * 61,
	})
	for i := range 63 {
		fmt.Printf("%d\n", i)
		timer.Tick()
	}
}
