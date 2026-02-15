//go:build linux

package xpad

import (
	"syscall"
	"time"
	"unsafe"
)

func waitReadable(fd int, timeout time.Duration) error {
	deadline := time.Time{}
	if timeout >= 0 {
		deadline = time.Now().Add(timeout)
	}

	for {
		var tv *syscall.Timeval
		if timeout >= 0 {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return ErrTimeout
			}
			t := syscall.NsecToTimeval(remaining.Nanoseconds())
			tv = &t
		}

		var readfds syscall.FdSet
		fdSet(fd, &readfds)

		n, err := syscall.Select(fd+1, &readfds, nil, nil, tv)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		if n == 0 {
			return ErrTimeout
		}
		return nil
	}
}

func fdSet(fd int, set *syscall.FdSet) {
	bitsPerWord := uint(unsafe.Sizeof(set.Bits[0]) * 8)
	idx := fd / int(bitsPerWord)
	if idx < 0 || idx >= len(set.Bits) {
		return
	}
	set.Bits[idx] |= 1 << (uint(fd) % bitsPerWord)
}
