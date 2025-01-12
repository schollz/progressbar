package main

import (
	"github.com/schollz/progressbar/v3"
	"time"
)

// This example illustrates a simple spinner that displays the number of
// received calls
func main() {
	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription("[Counter Example] Received"),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	c := make(chan int)

	go func(notify chan<- int) {
		defer close(notify)
		for i := 0; i < 100; i++ {
			c <- i
			time.Sleep(time.Millisecond * 100)
		}
	}(c)

	for v := range c {
		bar.Add(v)
	}
}
