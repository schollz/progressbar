package progressbar

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a thread-safe, simple
// progress bar
type ProgressBar struct {
	Max               int // max number of the counter
	Size              int // size of the saucer
	currentNum        int
	currentPercent    int
	lastPercent       int
	currentSaucerSize int

	lastShown time.Time
	startTime time.Time
	w         io.Writer

	// symbols
	symbolFinished string
	symbolLeft     string
	leftBookend    string
	rightBookend   string
	sync.RWMutex
}

// New returns a new ProgressBar
// with the specified maximum
func New(max int) *ProgressBar {
	p := new(ProgressBar)
	p.Lock()
	defer p.Unlock()
	p.Max = max
	p.Size = 40
	p.symbolFinished = "█"
	p.symbolLeft = " "
	p.leftBookend = "|"
	p.rightBookend = "|"
	p.w = os.Stdout
	p.lastShown = time.Now()
	p.startTime = time.Now()
	return p
}

// Add a certain amount to the progress bar
func (p *ProgressBar) Add(num int) error {
	p.Lock()
	p.currentNum += num
	percent := float64(p.currentNum) / float64(p.Max)
	p.currentSaucerSize = int(percent * float64(p.Size))
	p.currentPercent = int(percent * 100)
	updateBar := p.currentPercent != p.lastPercent
	p.lastPercent = p.currentPercent
	p.Unlock()
	if updateBar {
		return p.Show()
	}
	return nil
}

// Show will print the current progress bar
func (p *ProgressBar) Show() error {
	p.RLock()
	defer p.RUnlock()
	secondsLeft := time.Since(p.startTime).Seconds() / float64(p.currentNum) * (float64(p.Max) - float64(p.currentNum))
	s := fmt.Sprintf("\r%3d%% %s%s%s%s [%s:%s]",
		p.currentPercent,
		p.leftBookend,
		strings.Repeat(p.symbolFinished, p.currentSaucerSize),
		strings.Repeat(p.symbolLeft, p.Size-p.currentSaucerSize),
		p.rightBookend,
		time.Since(p.startTime).Round(time.Second).String(),
		(time.Duration(secondsLeft) * time.Second).String(),
	)

	_, err := io.WriteString(p.w, s)
	if err != nil {
		return err
	}
	if f, ok := p.w.(*os.File); ok {
		f.Sync()
	}
	return nil
}
