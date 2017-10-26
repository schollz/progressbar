# progressbar

[![travis](https://travis-ci.org/schollz/progressbar.svg?branch=master)](https://travis-ci.org/schollz/progressbar) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/progressbar)](https://goreportcard.com/report/github.com/schollz/progressbar) 
![coverage](https://img.shields.io/badge/coverage-94%25-brightgreen.svg)
[![godocs](https://godoc.org/github.com/schollz/progressbar?status.svg)](https://godoc.org/github.com/schollz/progressbar) 

A very simple progress bar.

**Basic usage:**

```golang
bar := progressbar.New(100)
for i := 0; i < 100; i++ {
    bar.Add(1)
    time.Sleep(10 * time.Millisecond)
}
```

which looks like:

```bash
 100% |████████████████████████████████████████| [1s:0s]            
 ```
