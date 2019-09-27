# progressbar

[![travis](https://travis-ci.org/schollz/progressbar.svg?branch=master)](https://travis-ci.org/schollz/progressbar) 
[![go report card](https://goreportcard.com/badge/github.com/schollz/progressbar)](https://goreportcard.com/report/github.com/schollz/progressbar) 
[![coverage](https://img.shields.io/badge/coverage-84%25-brightgreen.svg)](https://gocover.io/github.com/schollz/progressbar)
[![godocs](https://godoc.org/github.com/schollz/progressbar?status.svg)](https://godoc.org/github.com/schollz/progressbar) 

A very simple thread-safe progress bar which should work on every OS without problems. I needed a progressbar for [croc](https://github.com/schollz/croc) and everything I tried had problems, so I made another one. In order to be OS agnostic I do not plan to support [multi-line outputs](https://github.com/schollz/progressbar/issues/6).

![Example of progress bar](https://user-images.githubusercontent.com/6550035/32120326-5f420d42-bb15-11e7-89d4-c502864e78eb.gif)

## Install

```
go get -u github.com/schollz/progressbar/v2
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

### Progress for I/O operations

The `progressbar` implements an `io.Writer` so it can automatically detect the number of bytes written to a stream, so you can use it as a progressbar for an `io.Reader`.

```golang
urlToGet := "https://github.com/schollz/croc/releases/download/v4.1.4/croc_v4.1.4_Windows-64bit_GUI.zip"
req, _ := http.NewRequest("GET", urlToGet, nil)
resp, _ := http.DefaultClient.Do(req)
defer resp.Body.Close()

var out io.Writer
f, _ := os.OpenFile("croc_v4.1.4_Windows-64bit_GUI.zip", os.O_CREATE|os.O_WRONLY, 0644)
out = f
defer f.Close()

bar := progressbar.NewOptions(
    int(resp.ContentLength), 
    progressbar.OptionSetBytes(int(resp.ContentLength)),
)
out = io.MultiWriter(out, bar)
io.Copy(out, resp.Body)
```

See the tests for another example.

### Changing max value

The `progressbar` implements `ChangeMax` and `ChangeMax64` functions to change the max value of the progress bar.

```golang
bar := progressbar.New(100)
bar.ChangeMax(200) // Change the max of the progress bar to 200, not 100
```

You can also use `ChangeMax64` to minimize casting in the library.
See the tests for another example.

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

## License

MIT
