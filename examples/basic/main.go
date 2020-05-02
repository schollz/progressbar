package main

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	fmt.Println(getWidth())
	bar := progressbar.Default(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(20 * time.Millisecond)
	}
}
