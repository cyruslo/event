package event

import (
	"fmt"
	"runtime"
	"time"
)

func test(args ...interface{}) interface{} {
	return nil
}

func testEventMode() {
	initId := 1000000
	for i := 0; i < 1000; i++ {
		go Subscribe(1, NewHandleFunc(int64(initId+i), test))
	}
	fmt.Printf("goroutine num %d", runtime.NumGoroutine())
	time.Sleep(time.Second * 10)
	for i := 0; i < 1000; i++ {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Printf("panic %s", e)
				}
			}()
			Execute(1, int64(initId+i))
		}()
	}
	fmt.Printf("goroutine num %d", runtime.NumGoroutine())
}

func main() {
	go testEventMode()
}