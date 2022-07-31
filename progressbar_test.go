package progressbar

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRender(b *testing.B) {
	bar := NewOptions64(100000000,
		OptionSetWriter(os.Stderr),
		OptionShowIts(),
	)
	for i := 0; i < b.N; i++ {
		bar.Add(1)
	}
}

func ExampleProgressBar() {
	bar := New(100)
	bar.Add(10)
	// Output:
	// 10% |████                                    |  [0s:0s]
}

func ExampleProgressBar_Set() {
	bar := New(100)
	bar.Set(10)
	// Output:
	// 10% |████                                    |  [0s:0s]
}

func ExampleProgressBar_Set64() {
	bar := New(100)
	bar.Set64(10)
	// Output:
	// 10% |████                                    |  [0s:0s]
}

func ExampleProgressBar_basic() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         |  [1s:9s]
}

func ExampleProgressBar_invisible() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false), OptionSetVisibility(false))
	bar.Reset()
	bar.RenderBlank()
	fmt.Println("hello, world")
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// hello, world
}

func ExampleOptionThrottle() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false), OptionThrottle(100*time.Millisecond))
	bar.Reset()
	bar.Add(5)
	time.Sleep(150 * time.Millisecond)
	bar.Add(5)
	bar.Add(10)
	// Output:
	// 10% |█         |  [0s:1s]
}

func ExampleOptionClearOnFinish() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false), OptionClearOnFinish())
	bar.Reset()
	bar.Finish()
	fmt.Println("Finished")
	// Output:
	// Finished
}

func ExampleProgressBar_Finish() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false))
	bar.Finish()
	// Output:
	// 100% |██████████|
}

func Example_xOutOfY() {
	bar := NewOptions(100, OptionSetPredictTime(true))

	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(1 * time.Millisecond)
	}
}

func ExampleOptionShowIts_count() {
	bar := NewOptions(100, OptionSetWidth(10), OptionShowIts(), OptionShowCount())
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | (10/100, 10 it/s) [1s:9s]
}

func ExampleOptionShowIts() {
	bar := NewOptions(100, OptionSetWidth(10), OptionShowIts(), OptionSetPredictTime(false))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | (10 it/s)
}

func ExampleOptionShowCountBigNumber() {
	bar := NewOptions(10000, OptionSetWidth(10), OptionShowCount(), OptionSetPredictTime(false))
	bar.Add(1)
	// Output:
	// 0% |          | (1/10000)
}

func ExampleOptionSetPredictTime() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetPredictTime(false))
	_ = bar.Add(10)
	// Output:
	// 10% |█         |
}

func ExampleDefault() {
	bar := Default(100)
	for i := 0; i < 50; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
	// Output:
	//
}

func ExampleOptionChangeMax() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetPredictTime(false))
	bar.ChangeMax(50)
	bar.Add(50)
	// Output:
	// 100% |██████████|
}

func ExampleIgnoreLength_WithIteration() {
	/*
		IgnoreLength test with iteration count and iteration rate
	*/
	bar := NewOptions(-1,
		OptionSetWidth(10),
		OptionShowIts(),
		OptionShowCount(),
	)
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(5)

	// Output:
	// -  (5/-, 5 it/s)
}

func TestSpinnerType(t *testing.T) {
	bar := NewOptions(-1,
		OptionSetWidth(10),
		OptionSetDescription("indeterminate spinner"),
		OptionShowIts(),
		OptionShowCount(),
		OptionSpinnerType(9),
	)
	bar.Reset()
	for i := 0; i < 10; i++ {
		time.Sleep(120 * time.Millisecond)
		bar.Add(1)
	}
	if false {
		t.Errorf("error")
	}
}

func Test_IsFinished(t *testing.T) {
	isCalled := false
	bar := NewOptions(72, OptionOnCompletion(func() {
		isCalled = true
	}))

	// Test1: If bar is not fully completed.
	bar.Add(5)
	if bar.IsFinished() || isCalled {
		t.Errorf("Successfully tested bar is not yet finished.")
	}

	// Test2: Bar fully completed.
	bar.Add(67)
	if !bar.IsFinished() || !isCalled {
		t.Errorf("Successfully tested bar is finished.")
	}

	// Test3: If increases maximum bytes error should be thrown and
	// bar finished will remain false.
	bar.Reset()
	err := bar.Add(73)
	if err == nil || bar.IsFinished() {
		t.Errorf("Successfully got error when bytes increases max bytes, bar finished: %v", bar.IsFinished())
	}
}

func ExampleIgnoreLength_WithSpeed() {
	/*
		IgnoreLength test with iterations and count
	*/
	bar := NewOptions(-1,
		OptionSetWidth(10),
		OptionShowBytes(true),
	)

	bar.Reset()
	time.Sleep(1 * time.Second)
	// since 10 is the width and we don't know the max bytes
	// it will do a infinite scrolling.
	bar.Add(11)

	// Output:
	// -  (0.011 kB/s)
}

func TestBarSlowAdd(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions(100, OptionSetWidth(10), OptionShowIts(), OptionSetWriter(&buf))
	bar.Reset()
	time.Sleep(3 * time.Second)
	bar.Add(1)
	if !strings.Contains(buf.String(), "1%") {
		t.Errorf("wrong string: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "20 it/min") {
		t.Errorf("wrong string: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "[3s:") {
		t.Errorf("wrong string: %s", buf.String())
	}
	// Output:
	// 1% |          | (20 it/min) [3s:4m57s]
}

func TestBarSmallBytes(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions64(100000000, OptionShowBytes(true), OptionShowCount(), OptionSetWidth(10), OptionSetWriter(&buf))
	for i := 1; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Add(1000)
	}
	if !strings.Contains(buf.String(), "8.8 kB/95 MB") {
		t.Errorf("wrong string: %s", buf.String())
	}
	for i := 1; i < 10; i++ {
		time.Sleep(10 * time.Millisecond)
		bar.Add(1000000)
	}
	if !strings.Contains(buf.String(), "8.6/95 MB") {
		t.Errorf("wrong string: %s", buf.String())
	}
}

func TestBar(t *testing.T) {
	bar := New(0)
	if err := bar.Add(1); err == nil {
		t.Error("should have an error for 0 bar")
	}
	bar = New(10)
	if err := bar.Add(11); err == nil {
		t.Error("should have an error for adding > bar")
	}
}

func TestState(t *testing.T) {
	bar := NewOptions(100, OptionSetWidth(10))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	s := bar.State()
	if s.CurrentPercent != 0.1 {
		t.Error(s)
	}
}

func ExampleOptionSetRenderBlankState() {
	NewOptions(10, OptionSetWidth(10), OptionSetRenderBlankState(true))
	// Output:
	// 0% |          |  [0s:0s]
}

func TestBasicSets(t *testing.T) {
	b := NewOptions(
		999,
		OptionSetWidth(888),
		OptionSetRenderBlankState(true),

		OptionSetWriter(ioutil.Discard), // suppressing output for this test
	)

	tc := b.config

	if tc.max != 999 {
		t.Errorf("Expected %s to be %d, instead I got %d\n%+v", "max", 999, tc.max, b)
	}

	if tc.width != 888 {
		t.Errorf("Expected %s to be %d, instead I got %d\n%+v", "width", 999, tc.max, b)
	}

	if !tc.renderWithBlankState {
		t.Errorf("Expected %s to be %t, instead I got %t\n%+v", "renderWithBlankState", true, tc.renderWithBlankState, b)
	}
}

func TestOptionSetTheme(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions(
		10,
		OptionSetTheme(Theme{Saucer: "#", SaucerPadding: "-", BarStart: ">", BarEnd: "<"}),
		OptionSetWidth(10),
		OptionSetWriter(&buf),
	)
	bar.Add(5)
	result := strings.TrimSpace(buf.String())
	expect := "50% >#####-----<  [0s:0s]"
	if result != expect {
		t.Errorf("Render miss-match\nResult: '%s'\nExpect: '%s'\n%+v", result, expect, bar)
	}
}

// TestOptionSetPredictTime ensures that when predict time is turned off, the progress
// bar is showing the total steps completed of the given max, otherwise the predicted
// time in seconds is specified.
func TestOptionSetPredictTime(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions(
		10,
		OptionSetPredictTime(false),
		OptionSetWidth(10),
		OptionSetWriter(&buf),
	)

	_ = bar.Add(2)
	result := strings.TrimSpace(buf.String())
	expect := "20% |██        |"

	if result != expect {
		t.Errorf("Render miss-match\nResult: '%s'\nExpect: '%s'\n%+v", result, expect, bar)
	}

	bar.Reset()
	bar.config.predictTime = true
	buf.Reset()

	_ = bar.Add(7)
	result = strings.TrimSpace(buf.String())
	expect = "70% |███████   |  [0s:0s]"

	if result != expect {
		t.Errorf("Render miss-match\nResult: '%s'\nExpect: '%s'\n%+v", result, expect, bar)
	}
}

func TestIgnoreLength(t *testing.T) {
	bar := NewOptions(
		-1,
		OptionSetWidth(100),
	)
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)

	state := bar.State()
	if state.CurrentBytes != 10.0 {
		t.Errorf("Number of bytes mismatched gotBytes %f wantBytes %f", state.CurrentBytes, 10.0)
	}
	if state.CurrentPercent != 0.1 {
		t.Errorf("Percent of bar mismatched got %f want %f", state.CurrentPercent, 0.1)
	}

	kbPerSec := fmt.Sprintf("%2.2f", state.KBsPerSecond)
	if kbPerSec != "0.01" {
		t.Errorf("Speed mismatched got %s want %s", kbPerSec, "0.01")
	}
}

func TestReaderToBuffer(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	urlToGet := "https://dl.google.com/go/go1.14.1.src.tar.gz"
	req, err := http.NewRequest("GET", urlToGet, nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	bar := NewOptions(int(resp.ContentLength), OptionShowBytes(true), OptionShowCount())
	out := io.MultiWriter(buf, bar)
	_, err = io.Copy(out, resp.Body)
	assert.Nil(t, err)

	md5, err := md5sum(buf)
	assert.Nil(t, err)
	assert.Equal(t, "d441819a800f8c90825355dfbede7266", md5)
}

func TestReaderToFile(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	urlToGet := "https://dl.google.com/go/go1.14.1.src.tar.gz"
	req, err := http.NewRequest("GET", urlToGet, nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", "progressbar_testfile")
	if err != nil {
		t.Fatal()
	}
	defer os.Remove(f.Name())
	defer f.Close()

	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = fakeStdout

	bar := DefaultBytes(resp.ContentLength)
	out := io.MultiWriter(f, bar)
	_, err = io.Copy(out, resp.Body)
	assert.Nil(t, err)
	f.Sync()
	f.Seek(0, 0)

	if err := fakeStdout.Close(); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(b))

	md5, err := md5sum(f)
	assert.Nil(t, err)
	assert.Equal(t, "d441819a800f8c90825355dfbede7266", md5)
}

func TestReaderToFileUnknownLength(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	urlToGet := "https://dl.google.com/go/go1.14.1.src.tar.gz"
	req, err := http.NewRequest("GET", urlToGet, nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", "progressbar_testfile")
	if err != nil {
		t.Fatal()
	}
	defer os.Remove(f.Name())
	defer f.Close()

	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = fakeStdout

	bar := DefaultBytes(-1, " downloading")
	out := io.MultiWriter(f, bar)
	_, err = io.Copy(out, resp.Body)
	assert.Nil(t, err)
	f.Sync()
	f.Seek(0, 0)

	if err := fakeStdout.Close(); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(b))

	md5, err := md5sum(f)
	assert.Nil(t, err)
	assert.Equal(t, "d441819a800f8c90825355dfbede7266", md5)
}

func TestConcurrency(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions(
		1000,
		OptionSetWriter(&buf),
	)
	var wg sync.WaitGroup
	for i := 0; i < 900; i++ {
		wg.Add(1)
		go func(b *ProgressBar, wg *sync.WaitGroup) {
			bar.Add(1)
			wg.Done()
		}(bar, &wg)
	}
	wg.Wait()
	result := bar.state.currentNum
	expect := int64(900)
	assert.Equal(t, expect, result)
}

func TestIterationNames(t *testing.T) {
	b := Default(20)
	tc := b.config

	// Checking for the default iterations per second or "it/s"
	if tc.iterationString != "it" {
		t.Errorf("Expected %s to be %s, instead I got %s", "iterationString", "it", tc.iterationString)
	}

	// Change the default "it/s" to provide context, downloads per second or "dl/s"
	b = NewOptions(20, OptionSetItsString("dl"))
	tc = b.config

	if tc.iterationString != "dl" {
		t.Errorf("Expected %s to be %s, instead I got %s", "iterationString", "dl", tc.iterationString)
	}
}

func md5sum(r io.Reader) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, r)
	return hex.EncodeToString(hash.Sum(nil)), err
}

func TestProgressBar_Describe(t *testing.T) {
	buf := strings.Builder{}
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false), OptionSetWriter(&buf))
	bar.Describe("performing axial adjustments")
	bar.Add(10)
	rawBuf := strconv.QuoteToASCII(buf.String())
	if rawBuf != `"\rperforming axial adjustments   0% |          |  [0s:0s]\r                                                       \rperforming axial adjustments  10% |\u2588         |  [0s:0s]"` {
		t.Errorf("wrong string: %s", rawBuf)
	}
}
