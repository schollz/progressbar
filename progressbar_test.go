package progressbar

import (
	"testing"
	"time"
)

func ExampleProgressBar() {
	bar := NewOptions(OptionSetMax(100), OptionSetSize(10))
	bar.Reset()
	time.Sleep(1 * time.Second)
	bar.Add(10)
	// Output:
	// 10% |█         | [1s:9s]
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

func ExampleProgressBar_RenderBlank() {
	bar := NewOptions(OptionSetMax(10), OptionSetSize(10))
	bar.RenderBlank()
	// Output:
	// 0% |          | [0s:0s]
}

func TestSetMax(t *testing.T) {
	var b *ProgressBar
	expect := 999
	b = NewOptions(OptionSetMax(expect))

	if b.config.max != expect {
		t.Errorf("Expected max to be %d, instead I got %d\n%+v", expect, b.config.max, b)
	}

}
