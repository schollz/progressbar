package themes

import "testing"

func TestNew(t *testing.T) {
	if _, err := NewDefault(uint(len(defaultSymbolsFinished) + 1)); err == nil {
		t.Error("should have an error if n > default themes")
	}
}
