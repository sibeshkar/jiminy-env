package main

import (
	"fmt"
	"time"
)

func rewarder(ch chan int, d time.Duration) {
	var i int
	for {
		ch <- i
		i++
		time.Sleep(d)
	}
}

func envcontroller(ch chan int, d time.Duration) {
	var i int
	for {
		ch <- i
		i++
		time.Sleep(d)
	}
}

func reader(out chan int) {
	for x := range out {
		fmt.Println(x)
	}
}

// func main() {
// 	ch := make(chan int)
// 	out := make(chan int)
// 	go rewarder(ch, 100*time.Millisecond)
// 	go envcontroller(ch, 250*time.Millisecond)
// 	go reader(out)
// 	for i := range ch {
// 		out <- i
// 	}
// }
