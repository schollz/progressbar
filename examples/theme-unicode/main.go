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
		progressbar.OptionSetTheme(progressbar.ThemeUnicode),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionFullWidth(),
	)

	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(66 * time.Millisecond)
	}
}
