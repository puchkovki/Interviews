// echoloop
package main

import (
	"fmt"
	"os"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func echo(ch1 chan []string, quit chan bool) {
	firstborn := make([]string, 0)
	for {
		select {
		case value := <-ch1:
			firstborn = append(firstborn, value...)
		case <-quit:
			fmt.Println("quit")
			return
			//return &MyError{time.Now(), "You cheated on me"}
		default:
			fmt.Println(firstborn)
			time.Sleep(time.Second)
		}
	}
	return
}

func main() {
	argv := os.Args[1:]
	ch1 := make(chan []string)
	quit := make(chan bool)
	go echo(ch1, quit)
	ch1 <- argv
}
