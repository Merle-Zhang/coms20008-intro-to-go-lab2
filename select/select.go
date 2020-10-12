package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

// slowSender sends a string every 2 seconds.
func slowSender(c chan<- string) {
	for {
		time.Sleep(2 * time.Second)
		c <- "I am the slowSender"
	}
}

// fastSender sends consecutive ints every 500 ms.
func fastSender(c chan<- int) {
	for i := 0; ; i++ {
		time.Sleep(500 * time.Millisecond)
		c <- i
	}
}

func fasterSender(c chan<- []int) {
	slice := []int{1, 2, 3}
	for {
		time.Sleep(200 * time.Millisecond)
		c <- slice
	}
}

func helper() {
	ints := make(chan int, 10)
	go fastSender(ints)
	strings := make(chan string, 10)
	go slowSender(strings)
	slices := make(chan []int, 10)
	go fasterSender(slices)

	for { // = while(true)
		select {
		case s := <-strings:
			fmt.Println("Received a string", s)
		case i := <-ints:
			fmt.Println("Received an int", i)
		case sl := <-slices:
			fmt.Println("Received a slice", sl)
		default:
			fmt.Println("--- Nothing to receive, sleeping for 3s...")
			time.Sleep(3 * time.Second)
		}
	}
}

// main starts the two senders and then goes into an infinite loop of receiving their messages.
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

	go helper()
	time.Sleep(7 * time.Second)
}
