package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	// TODO: Write an infinite loop of sending "pings" and receiving "pongs"
	message := ""
	for {
		message = "ping"
		channel <- message
		fmt.Println("Foo is sending:", message)
		// time.Sleep(500 * time.Millisecond) // TODO: why don't need this line?
		// Ans: If the channel isn't buffered, the program will sleep when it sends until the channel is read from, so it isn't possible to receive what you just sent
		message = <-channel
		fmt.Println("Foo has received:", message)
		fmt.Println()
	}
}

func bar(channel chan string) {
	// TODO: Write an infinite loop of receiving "pings" and sending "pongs"
	message := ""
	for {
		message = <-channel
		fmt.Println("Bar has received:", message)
		message = "pong"
		channel <- message
		fmt.Println("Bar is sending:", message)
	}
}

func pingPong() {
	// TODO: make channel of type string and pass it to foo and bar
	channel := make(chan string)
	go foo(channel) // Nil is similar to null. Sending or receiving from a nil chan blocks forever.
	go bar(channel)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
