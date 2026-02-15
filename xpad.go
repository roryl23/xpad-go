package xpad

import (
	"errors"
	"os"
	"time"
)

var (
	ErrClosed         = errors.New("xpad: device is closed")
	ErrNotImplemented = errors.New("xpad: not implemented")
	ErrNotFound       = errors.New("xpad: no matching device found")
	ErrReadOnly       = errors.New("xpad: device opened read-only")
	ErrTimeout        = errors.New("xpad: read timeout")
)

// Device represents an open xpad device.
type Device struct {
	Path     string
	file     *os.File
	readOnly bool
}

// Event represents an input_event from the Linux input subsystem.
type Event struct {
	When  time.Time
	Kind  EventKind
	Code  uint16
	Value int32
}

// EventKind classifies an input event (EV_*).
type EventKind uint16

const (
	EventUnknown EventKind = 0
)

// Open opens an xpad device by path.
func Open(path string) (*Device, error) {
	file, readOnly, err := openReadWriteOrReadOnly(path)
	if err != nil {
		return nil, err
	}
	return &Device{Path: path, file: file, readOnly: readOnly}, nil
}

// Close closes the device.
func (d *Device) Close() error {
	if d == nil || d.file == nil {
		return nil
	}
	err := d.file.Close()
	d.file = nil
	d.readOnly = false
	return err
}

// File returns the underlying file handle when open.
func (d *Device) File() (*os.File, error) {
	if d == nil || d.file == nil {
		return nil, ErrClosed
	}
	return d.file, nil
}

// FD returns the underlying file descriptor when open.
func (d *Device) FD() (uintptr, error) {
	if d == nil || d.file == nil {
		return 0, ErrClosed
	}
	return d.file.Fd(), nil
}

// ReadOnly reports whether the device handle is opened read-only.
func (d *Device) ReadOnly() bool {
	if d == nil {
		return true
	}
	return d.readOnly
}

// OpenDevice opens an xpad device from discovery info.
func OpenDevice(info DeviceInfo) (*Device, error) {
	return Open(info.Path)
}

// ReadEvent blocks until the next event or timeout.
func (d *Device) ReadEvent(timeout time.Duration) (Event, error) {
	if d == nil || d.file == nil {
		return Event{}, ErrClosed
	}
	return readEvent(d, timeout)
}

// SendEvent writes an input_event to the device.
func (d *Device) SendEvent(ev Event) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	return writeEvent(d, ev)
}
