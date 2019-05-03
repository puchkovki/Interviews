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
	file.Close()
	f, _ := file.Stat()
	os.Remove(f.Name())

	fifo.Close()
	f, _ = fifo.Stat()
	os.Remove(f.Name())

	os.Exit(0)
}

func echo(channel chan []string) {
	firstborn := make([]string, 0)
	for {
		select {
		case value := <-channel:
			firstborn = append(firstborn, value...)
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
		writer := bufio.NewWriter(fifo)
		for _, value := range argv {
			writer.WriteString(value)
		}
		writer.Flush()
		return
	} else {
		fmt.Println("Successfully locked the file")
	}

	os.Remove(pipe)
	err = syscall.Mkfifo(pipe, 0600)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	} else {
		fmt.Println("Successfully made named pipe")
	}
	fifo, err := os.OpenFile(pipe, os.O_RDWR|os.O_CREATE, os.ModeNamedPipe) //ModeNamedPipe = named pipe(fifo)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	} else {
		fmt.Println("Successfully opened named pipe")
	}
	reader := bufio.NewReader(fifo)

	channel := make(chan []string)

	sigs := make(chan os.Signal, 1) //signal handler
	signal.Notify(sigs, os.Kill, os.Interrupt)
	//Эта горутина будет заблокирована, пока мы не получим сигнал.
	//При его получении, он будет выведен и мы оповестим программу о том, что она может прекратить свое выполнение.
	go Cleaner(sigs, fifo, file)

	go echo(channel)
	channel <- argv
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
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
