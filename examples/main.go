package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v2"
)

func main() {
	fmt.Println("downloading go1.12.5.linux-amd64.tar.gz")
	defer os.Remove("go1.12.5.linux-amd64.tar.gz")
	urlToGet := "https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz"
	req, _ := http.NewRequest("GET", urlToGet, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var out io.Writer
	f, _ := os.OpenFile("go1.12.5.linux-amd64.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	out = f
	defer f.Close()

	bar := progressbar.NewOptions(
		int(resp.ContentLength),
		progressbar.OptionSetBytes(int(resp.ContentLength)),
		progressbar.OptionThrottle(10*time.Millisecond),
	)
	out = io.MultiWriter(out, bar)
	io.Copy(out, resp.Body)
	fmt.Println("done")

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

	bar = progressbar.NewOptions(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetBytes(10000),
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
