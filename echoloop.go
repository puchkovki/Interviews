// echoloop
package main

import (
	//"bufio"
	"fmt"
	//"io"
	"log"
	"os"
	"syscall"
	"time"
)

func echo(channel chan []string, quit chan bool) {
	firstborn := make([]string, 0)
	for {
		select {
		case value := <-channel:
			firstborn = append(firstborn, value...)
		case <-quit:
			fmt.Println("quit")
			return
		default:
			fmt.Println(firstborn, "Okay")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	argv := os.Args[1:]
	pipe := "Echoloop.pipe"

	file, err := os.OpenFile("echoloop.lock", os.O_RDWR|os.O_CREATE, 0600) //read, write, not execute
	if err != nil {
		log.Fatal("Cannot open the file!", err)
	}
	//defer file.Close() useless piece of ...

	err = syscall.FcntlFlock(file.Fd(), syscall.F_SETLK, &syscall.Flock_t{Type: syscall.F_WRLCK, Pid: int32(os.Getpid())})
	if err != nil {
		fmt.Println("Cannot lock the file!")
		fifo, err := os.OpenFile(pipe, os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			log.Fatal("Open named pipe file error:", err)
		}
		for _, value := range argv {
			fifo.WriteString(value)
		}
		return
	}

	os.Remove(pipe)
	err = syscall.Mkfifo(pipe, 0600)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}
	fifo, err := os.OpenFile(pipe, os.O_CREATE, os.ModeNamedPipe) //ModeNamedPipe = named pipe(fifo)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}
	reader := bufio.NewReader(fifo)

	channel := make(chan []string)
	quit := make(chan bool) //signal handler
	//Handler(quit)
	go echo(channel, quit)
	channel <- argv
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		str := make([]string, 1)
		str = append(str, line)
		channel <- str
	}
}

/*func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error{}
syscall.F_SETLK acquire a lock
parameters of F_SETLK: type Flock_t struct {Pid int32 — Getpid returns int*/
