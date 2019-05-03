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

//Cleaner closes all required files and exits the program
func Cleaner(sigs chan os.Signal, mutex, fifo *os.File) {
	//Signal awaiting mode
	if sigs != nil {
		<-sigs
	}

	//Deleting the mutex file
	if mutex != nil {
		f, _ := mutex.Stat()
		mutex.Close()
		os.Remove(f.Name())
	} else {
		fmt.Println("mutex file == nil")
		return
	}

	//Deleting the fifo file
	if fifo != nil {
		f, _ := fifo.Stat()
		fifo.Close()
		os.Remove(f.Name())
	} else {
		fmt.Println("fifo file == nil")
		return
	}

	fmt.Println()
	os.Exit(0)
}

//echo prints receined arguments
func echo(channel chan []string) {
	firstborn := make([]string, 0)

	for {
		select {
		case value := <-channel:
			firstborn = append(firstborn, value...)
		default:
			for _, value := range firstborn {
				fmt.Println(value)
			}
			time.Sleep(time.Second)
		}
	}
}

func main() {
	argv := os.Args[1:]
	pipeName := "Echoloop.pipe"

	//Opening file, if it doesn't exist — create one
	mutex, err := os.OpenFile("echoloop.lock", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		Cleaner(nil, mutex, nil)
		log.Fatal("Cannot open the mutex file!", err)
	}

	//Try to lock the mutex
	err = syscall.FcntlFlock(mutex.Fd(), syscall.F_SETLK, &syscall.Flock_t{Type: syscall.F_WRLCK, Pid: int32(os.Getpid())})

	//Not the firstborn
	if err != nil {
		fmt.Println("Cannot lock the mutex file!")

		//Putting arguments into the fifo file
		fifo, err := os.OpenFile(pipeName, os.O_RDWR|os.O_CREATE, os.ModeNamedPipe)
		if err != nil {
			Cleaner(nil, mutex, fifo)
			log.Fatal("Open named pipe file error:", err)
		}

		//Making a buffer to send arguments to the firstborn
		writer := bufio.NewWriter(fifo)
		for _, value := range argv {
			writer.WriteString(value)
		}

		//Marking the end of the transfer
		writer.WriteByte(0)
		writer.Flush()
		return
	}

	//Deleting the fifo file of it already exists
	os.Remove(pipeName)

	//Creating a fifo file
	err = syscall.Mkfifo(pipeName, 0600)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	//Opening fifo file
	fifo, err := os.OpenFile(pipeName, os.O_RDWR|os.O_CREATE, os.ModeNamedPipe) //ModeNamedPipe = named pipe(fifo)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	//Making a buffer to get arguments from bastards
	reader := bufio.NewReader(fifo)

	//Making a channel for sending and receiving strings to print
	channel := make(chan []string)

	//Making a channel for handling signals
	sigs := make(chan os.Signal, 1)

	//Catching os.Kill and os.Interrupt signals
	signal.Notify(sigs, os.Kill, os.Interrupt)
	go Cleaner(sigs, mutex, fifo)

	//Printing
	go echo(channel)
	channel <- argv

	//Receiving arguments from the fifo file
	for {
		line, err := reader.ReadString(0)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				Cleaner(sigs, mutex, fifo)
				log.Fatal(err)
			}
		}

		fmt.Println(line)

		str := make([]string, 1)
		str[0] = line
		channel <- str
	}
}
