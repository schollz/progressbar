package progressbar

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleProgressBar() {
	bar := New(100)
	bar.Add(10)
	// Output:
	// 10% |████                                    |  [0s:0s]
}
func ExampleProgressBarBasic() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         |  [1s:9s]
}

func ExampleThrottle() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false), OptionThrottle(100*time.Millisecond))
	bar.Reset()
	bar.Add(5)
	time.Sleep(150 * time.Millisecond)
	bar.Add(5)
	bar.Add(10)
	// Output:
	// 10% |█         |  [0s:1s]
}
func ExampleFinish() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false))
	bar.Reset()
	bar.Finish()
	// Output:
	// 100% |██████████|  [0s:0s]
}

func ExampleSetBytes() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetBytes(10000))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | (1.0 kB/s) [1s:9s]
}

func ExampleShowCount() {
	bar := NewOptions(100, OptionSetWidth(10), OptionShowIts(), OptionShowCount())
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | (10/100, 10 it/s) [1s:9s]
}

func ExampleSetIts() {
	bar := NewOptions(100, OptionSetWidth(10), OptionShowIts())
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | (10 it/s) [1s:9s]
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
	bar := NewOptions(100, OptionSetWidth(10), OptionSetBytes(10000))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	s := bar.State()
	if s.CurrentPercent != 0.1 {
		t.Error(s)
	}
}

func ExampleProgressBar_RenderBlank() {
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

func TestReaderToBuffer(t *testing.T) {
	urlToGet := "https://github.com/schollz/croc/releases/download/v4.1.4/croc_v4.1.4_Windows-64bit_GUI.zip"
	req, err := http.NewRequest("GET", urlToGet, nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	var out io.Writer
	/// setup buffer
	var buf bytes.Buffer
	f := bufio.NewWriter(&buf)
	out = f

	bar := NewOptions(int(resp.ContentLength), OptionSetBytes(int(resp.ContentLength)))
	out = io.MultiWriter(out, bar)
	_, err = io.Copy(out, resp.Body)
	assert.Nil(t, err)

	// if reading to buffer, write buffer bytes
	f.Flush()
	err = ioutil.WriteFile("croc_v4.1.4_Windows-64bit_GUI.zip", buf.Bytes(), 0644)
	assert.Nil(t, err)

	md5, err := md5sum("croc_v4.1.4_Windows-64bit_GUI.zip")
	assert.Nil(t, err)
	assert.Equal(t, "1e496ef2beba6e2a5e4200cba72a5ad6", md5)
	assert.Nil(t, os.Remove("croc_v4.1.4_Windows-64bit_GUI.zip"))
}

func TestReaderToFile(t *testing.T) {
	urlToGet := "https://github.com/schollz/croc/releases/download/v4.1.4/croc_v4.1.4_Windows-64bit_GUI.zip"
	req, err := http.NewRequest("GET", urlToGet, nil)
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	var out io.Writer
	// read to file
	f, err := os.OpenFile("croc_v4.1.4_Windows-64bit_GUI.zip", os.O_CREATE|os.O_WRONLY, 0666)
	assert.Nil(t, err)
	out = f

	bar := NewOptions(int(resp.ContentLength), OptionSetBytes(int(resp.ContentLength)))
	out = io.MultiWriter(out, bar)
	_, err = io.Copy(out, resp.Body)
	assert.Nil(t, err)
	f.Close()

	md5, err := md5sum("croc_v4.1.4_Windows-64bit_GUI.zip")
	assert.Nil(t, err)
	assert.Equal(t, "1e496ef2beba6e2a5e4200cba72a5ad6", md5)
	assert.Nil(t, os.Remove("croc_v4.1.4_Windows-64bit_GUI.zip"))
}

func md5sum(filePath string) (result string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func ExampleDescribe() {
	bar := NewOptions(100, OptionSetWidth(10), OptionSetRenderBlankState(false))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Describe("performing axial adjustements")
	bar.Add(10)
	// Output:
	// performing axial adjustements  10% |█         |  [1s:9s]
}
