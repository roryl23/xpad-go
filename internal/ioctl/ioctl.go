//go:build linux

// Package ioctl provides helpers for constructing and issuing Linux ioctl calls.
package ioctl

import (
	"syscall"
	"unsafe"
)

const (
	iocNone  = 0
	iocWrite = 1
	iocRead  = 2
)

// Direction flags for ioctl calls.
const (
	DirNone      = iocNone
	DirWrite     = iocWrite
	DirRead      = iocRead
	DirReadWrite = iocRead | iocWrite
)

const (
	iocNRBits   = 8
	iocTypeBits = 8
	iocSizeBits = 14
	iocDirBits  = 2
)

const (
	iocNRShift   = 0
	iocTypeShift = iocNRShift + iocNRBits
	iocSizeShift = iocTypeShift + iocTypeBits
	iocDirShift  = iocSizeShift + iocSizeBits
)

// IOC builds an ioctl request value.
func IOC(dir, typ, nr, size uint) uint {
	return (dir << iocDirShift) | (typ << iocTypeShift) | (nr << iocNRShift) | (size << iocSizeShift)
}

// IO builds an ioctl request with no data transfer.
func IO(typ, nr uint) uint {
	return IOC(iocNone, typ, nr, 0)
}

// IOR builds an ioctl request that reads data from the kernel.
func IOR(typ, nr, size uint) uint {
	return IOC(iocRead, typ, nr, size)
}

// IOW builds an ioctl request that writes data to the kernel.
func IOW(typ, nr, size uint) uint {
	return IOC(iocWrite, typ, nr, size)
}

// IOWR builds an ioctl request that reads and writes data.
func IOWR(typ, nr, size uint) uint {
	return IOC(iocRead|iocWrite, typ, nr, size)
}

// Size returns the size of a value in bytes for ioctl sizing.
func Size(v any) uint {
	return uint(unsafe.Sizeof(v))
}

// Call issues an ioctl with a raw argument pointer.
func Call(fd uintptr, req uint, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(req), arg)
	if errno != 0 {
		return errno
	}
	return nil
}

// CallPtr issues an ioctl using a typed pointer.
func CallPtr(fd uintptr, req uint, ptr unsafe.Pointer) error {
	return Call(fd, req, uintptr(ptr))
}
