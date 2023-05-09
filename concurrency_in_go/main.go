package main

import (
	"fmt"
	"sync"
	"time"
)

func writer(ch chan<- int, lock *sync.Mutex) {
	for i := 0; i < 100; i++ {
		lock.Lock()
		fmt.Println("Writing")
		ch <- i
		lock.Unlock()
	}
	close(ch)
}
func reader(ch chan int) {
	for i := range ch {
		fmt.Println("Reading ", i)
	}
}
func main() {
	fmt.Println("starting the program")
	ch := make(chan int)
	lock := &sync.Mutex{}
	go writer(ch, lock)
	go reader(ch)
	time.Sleep(5 * time.Second)
}
