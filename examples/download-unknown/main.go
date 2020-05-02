package main

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func main() {
	req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile("go1.14.2.src.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		-1,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
