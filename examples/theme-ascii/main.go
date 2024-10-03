package main

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	bar := progressbar.NewOptions(100,
		progressbar.OptionUseANSICodes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.ThemeASCII),
		progressbar.OptionShowElapsedTimeOnFinish(),
	)

	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(40 * time.Millisecond)
	}
}
