// +build linux darwin freebsd

package progressbar

import (
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() int {
	ws := &winsize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	width := int(uint(ws.Col))
	if width < 1 || width > 200 {
		return 80
	} else {
		return width
	}
}
