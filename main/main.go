package main

import (
	"github.com/MrMe42/progressbar"
//	"time"
)

func main() {
	bar := progressbar.New(100)
	bar.Add(50)
}

