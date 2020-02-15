package main

import (
	"fmt"
	"time"
)

var i = 0

func main() {
	a := make(chan bool)

	go func() {
		select {
		case <-a:
			fmt.Println(`12345`)
			break
		}
	}()

	go func() {
		select {
		case <-a:
			fmt.Println(`123`)
			break
		}
	}()

	close(a)

	time.Sleep(6 * time.Second)
}