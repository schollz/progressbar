package main

import (
	"time"

	"github.com/schollz/progressbar"
)

func main() {
	bar := progressbar.New(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
}
