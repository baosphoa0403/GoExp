package concurrency

import (
	"fmt"
	"time"
)

func Boring(msg string) {
	for i := 0; ; i++ {
		fmt.Printf("%s %d\n", msg, i)
		time.Sleep(time.Second)
	}
}

func MainBoring() {
	go Boring("boring!")
	fmt.Println("I'm listening.")
	time.Sleep(10 * time.Second)
	fmt.Println("You're boring; I'm leaving.")
}
