package main

import (
	"fmt"
	"log"
	"time"
)

func display() {
	const tickRate = 2 * time.Second
	stopper := time.After(10 * time.Second)
	ticker := time.NewTicker(tickRate).C
	log.Println("start")
loop:
	for {
		select {
		case <-ticker:
			fmt.Println("tick")
		case <-stopper:
			break loop
		}

	}

	log.Println("finish")
}
