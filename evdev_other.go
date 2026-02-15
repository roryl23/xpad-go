//go:build !linux

package xpad

import "time"

// Name is not supported on non-Linux platforms.
func (d *Device) Name() (string, error) { return "", ErrNotImplemented }

// Phys is not supported on non-Linux platforms.
func (d *Device) Phys() (string, error) { return "", ErrNotImplemented }

// Uniq is not supported on non-Linux platforms.
func (d *Device) Uniq() (string, error) { return "", ErrNotImplemented }

// ID is not supported on non-Linux platforms.
func (d *Device) ID() (InputID, error) { return InputID{}, ErrNotImplemented }

// AbsInfo is not supported on non-Linux platforms.
func (d *Device) AbsInfo(code uint16) (AbsInfo, error) { return AbsInfo{}, ErrNotImplemented }

// EventTypes is not supported on non-Linux platforms.
func (d *Device) EventTypes() ([]byte, error) { return nil, ErrNotImplemented }

// HasEventType is not supported on non-Linux platforms.
func (d *Device) HasEventType(ev EventKind) (bool, error) { return false, ErrNotImplemented }

// HasEventCode is not supported on non-Linux platforms.
func (d *Device) HasEventCode(ev EventKind, code uint16) (bool, error) {
	return false, ErrNotImplemented
}

// EffectCount is not supported on non-Linux platforms.
func (d *Device) EffectCount() (int, error) { return 0, ErrNotImplemented }

// Grab is not supported on non-Linux platforms.
func (d *Device) Grab(grab bool) error { return ErrNotImplemented }

func readEvent(d *Device, timeout time.Duration) (Event, error) {
	return Event{}, ErrNotImplemented
}

func writeEvent(d *Device, ev Event) error {
	return ErrNotImplemented
}
