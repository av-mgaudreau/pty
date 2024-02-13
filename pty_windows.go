package pty

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"unsafe"
)

// TODO: Dynamically find where cygwin is located
const cygwinPath = `C:\cygwin64`

func open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile(filepath.Join(cygwinPath, "dev", "ptmx"), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, err
	}
	// In case of error after this point, make sure we close the ptmx fd.
	defer func() {
		if err != nil {
			_ = p.Close() // Best effort.
		}
	}()

	// Make sure the pts folder exists, as Windows won't automatically create it
	os.Mkdir(filepath.Join(cygwinPath, "dev", "pts"), 0666)

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	if err := unlockpt(p); err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR|os.O_CREATE|syscall.O_NOCTTY, 0666) //nolint:gosec // Expected Open from a variable.
	if err != nil {
		return nil, nil, err
	}
	return p, t, nil
}

// TIOCGPTN                         = 0x80045430
func ptsname(f *os.File) (string, error) {
	var n _C_uint
	err := ioctl(f, 0x80045430, uintptr(unsafe.Pointer(&n))) //nolint:gosec // Expected unsafe pointer for Syscall call.
	if err != nil {
		return "", err
	}
	return filepath.Join(cygwinPath, "dev", "pts", strconv.Itoa(int(n))), nil
}

// TIOCSPTLCK = 0x40045431
func unlockpt(f *os.File) error {
	var u _C_int
	// use TIOCSPTLCK with a pointer to zero to clear the lock.
	return ioctl(f, 0x40045431, uintptr(unsafe.Pointer(&u))) //nolint:gosec // Expected unsafe pointer for Syscall call.
}
