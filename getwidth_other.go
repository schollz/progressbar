// +build !linux
// +build !darwin
// +build !freebsd
// +build !nacl

package progressbar

func getWidth() int {
	return 80
}
