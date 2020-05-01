# progressbar

[![travis](https://travis-ci.org/schollz/progressbar.svg?branch=master)](https://travis-ci.org/schollz/progressbar) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/progressbar)](https://goreportcard.com/report/github.com/schollz/progressbar) 
[![coverage](https://img.shields.io/badge/coverage-84%25-brightgreen.svg)](https://gocover.io/github.com/schollz/progressbar)
[![godocs](https://godoc.org/github.com/schollz/progressbar?status.svg)](https://godoc.org/github.com/schollz/progressbar) 

A very simple thread-safe progress bar which should work on every OS without problems. I needed a progressbar for [croc](https://github.com/schollz/croc) and everything I tried had problems, so I made another one. In order to be OS agnostic I do not plan to support [multi-line outputs](https://github.com/schollz/progressbar/issues/6).

![Example of progress bar](https://user-images.githubusercontent.com/6550035/32120326-5f420d42-bb15-11e7-89d4-c502864e78eb.gif)

## Install

```
go get -u github.com/schollz/progressbar/v3
```

## Usage 

### Basic usage

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

The times at the end show the elapsed time and the remaining time, respectively.

## Progress bar with unknown length

A progressbar with unknown length is a spinner. You can pass `-1` as the iterations to the bar and it will automatically convert it to a spinner with a customizable spinner type.

```golang
// basic bar
bar := progressbar.NewOptions(-1,
    progressbar.OptionSetDescription("indeterminate bar"),
    progressbar.OptionSpinnerType(70),
    progressbar.OptionSetWriter(os.Stderr),
    progressbar.OptionShowIts(),
    progressbar.OptionShowCount(),
)
for i := 0; i < 7000; i++ {
    bar.Add(1)
    time.Sleep(2 * time.Millisecond)
}
```

which looks like:

```bash
\ indeterminate bar (3969/-, 916 it/s)
```


### Progress for I/O operations

The `progressbar` implements an `io.Writer` so it can automatically detect the number of bytes written to a stream, so you can use it as a progressbar for an `io.Reader`.

```golang
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
    progressbar.OptionThrottle(10*time.Millisecond),
    progressbar.OptionShowCount(),
    progressbar.OptionOnCompletion(func() {
        fmt.Println(" done.")
    }),
)
io.Copy(io.MultiWriter(f, bar), resp.Body)
```

which looks like:

```bash
https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz 100% |██████████| (128/128 MB, 103.751 MB/s) [1s:0s] done.
```


### Long running processes

For long running processes, you might want to render from a 0% state.

```golang
// Renders the bar right on construction
bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
```

Alternatively, when you want to delay rendering, but still want to render a 0% state
```golang
bar := progressbar.NewOptions(100)

// Render the current state, which is 0% in this case
bar.RenderBlank()

// Emulate work
for i := 0; i < 10; i++ {
    time.Sleep(10 * time.Minute)
    bar.Add(10)
}
```

### Use a custom writer

The default writer is standard output (os.Stdout), but you can set it to whatever satisfies io.Writer.
```golang
bar := NewOptions(
    10,
    OptionSetTheme(Theme{Saucer: "#", SaucerPadding: "-", BarStart: ">", BarEnd: "<"}),
    OptionSetWidth(10),
    OptionSetWriter(&buf),
)

bar.Add(5)
result := strings.TrimSpace(buf.String())

// Result equals:
// 50% >#####-----< [0s:0s]

```



### Displaying Total Increment Over Predicted Time

By default the progress bar will attempt to predict the remaining amount of time left. This can be change to 
just show the current increment over the total maximum amount set for the progress bar. Do this by using the
`OptionSetPredictTime` option during progress bar creation.

```golang
bar := progressbar.NewOptions(100, progressbar.OptionSetPredictTime(false))
bar.Add(20)

// this result equals:
// "20% |██        |  [20:100]"

// default result equals:
// "20% |██        |  [3s:15s]"
```

## Contributing

Pull requests are welcome. Feel free to...

- Revise documentation
- Add new features
- Fix bugs
- Suggest improvements

## Thanks

Thanks [@Dynom](https://github.com/dynom) for massive improvements in version 2.0!

Thanks [@CrushedPixel](https://github.com/CrushedPixel) for adding descriptions and color code support!

Thanks [@MrMe42](https://github.com/MrMe42) for adding some minor features!

Thanks [@tehstun](https://github.com/tehstun) for some great PRs!

Thanks [@Benzammour](https://github.com/Benzammour) and [@haseth](https://github.com/haseth) for helping create v3!

## License

MIT
