package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

func main() {
	defer os.Remove("go1.12.5.linux-amd64.tar.gz")

	urlToGet := "https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz"
	req, _ := http.NewRequest("GET", urlToGet, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.12.5.linux-amd64.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionSetDescription(urlToGet),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(10*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println(" done.")
		}),
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	// basic bar
	bar = progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("indeterminate bar"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(70),
	)
	bar.RenderBlank() // will show the progress bar
	bar.Add(3000)
	time.Sleep(1 * time.Second)
	for i := 0; i < 7000; i++ {
		bar.Add(1)
		time.Sleep(2 * time.Millisecond)
	}

	// basic bar
	bar = progressbar.NewOptions(10000,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "~", SaucerPadding: "-", BarStart: "|", BarEnd: "|"}),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
	)
	bar.RenderBlank() // will show the progress bar
	bar.Add(3000)
	time.Sleep(1 * time.Second)
	for i := 0; i < 7000; i++ {
		bar.Add(1)
		time.Sleep(2 * time.Millisecond)
	}
	fmt.Println("finished1")

	// bar with options
	bar = progressbar.NewOptions(1000,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "~", SaucerPadding: "-", BarStart: "|", BarEnd: "|"}),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowIts(),
	)
	for i := 0; i < 1000; i++ {
		bar.Add(1)
		time.Sleep(2 * time.Millisecond)
	}
	fmt.Println("finished2")

	bar = progressbar.NewOptions(1000,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for i := 0; i < 1000; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("finished3")
}
