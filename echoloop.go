// echoloop
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Cleaner(sigs chan os.Signal, file, fifo *os.File) {
	<-sigs
	os.Close(file)
	name := file.Stat().Name()
	os.Remove(name)

	os.Close(fifo)
	name = fifo.Stat().Name()
	os.Remove(name)
	
	exit 0
}

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
	} else {
		fmt.Println("Successfully open the file")
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
	} else {
		fmt.Println("Successfullt locked the file")
	}

	os.Remove(pipe)
	err = syscall.Mkfifo(pipe, 0600)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	} else {
		fmt.Println("Successfullt made named pipe")
	}
	fifo, err := os.OpenFile(pipe, os.O_CREATE, os.ModeNamedPipe) //ModeNamedPipe = named pipe(fifo)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	} else {
		fmt.Println("Successfullt opened named pipe")
	}
	reader := bufio.NewReader(fifo)

	channel := make(chan []string)

	quit := make(chan bool, 1) //signal handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Kill, os.Interrupt)
	//Эта горутина будет заблокирована, пока мы не получим сигнал.
	//При его получении, он будет выведен и мы оповестим программу о том, что она может прекратить свое выполнение.
	go Cleaner(sigs, fifo, file)

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
