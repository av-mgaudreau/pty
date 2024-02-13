//go:build !solaris && !aix
// +build !solaris,!aix

package pty

// Local syscall const values.
const (
// TIOCGWINSZ = syscall.TIOCGWINSZ
// TIOCSWINSZ = syscall.TIOCSWINSZ
)

// SYS_IOCTL                  = 16
func ioctlInner(fd, cmd, ptr uintptr) error {
	// _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, ptr)
	// if e != 0 {
	// 	return e
	// }
	return nil
}
