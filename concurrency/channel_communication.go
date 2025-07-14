package concurrency

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func ChannelBoring(msg string, c chan string, cancel context.CancelFunc) {
	for i := 0; i < 5; i++ {
		c <- fmt.Sprintf("%s %d", msg, i)
		if i == 3 {
			cancel()
		}
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func MainChannelBoring() {
	c := make(chan string, 2)
	context, cancel := context.WithCancel(context.Background())
	go ChannelBoring("boring!", c, cancel)
	// for i := 0; i < 20; i++ {
	// 	fmt.Printf("You say: %q\n", <-c)
	// }

	for {
		select {
		case <-context.Done():
			fmt.Println("Cancel You're boring; I'm leaving.")
			return
		case res := <-c:
			fmt.Printf("You say: %q\n", res)
		}
	}
}
