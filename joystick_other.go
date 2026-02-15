//go:build !linux

package xpad

import "time"

// Joystick represents an open /dev/input/js* device.
type Joystick struct{}

// OpenJoystick is not supported on non-Linux platforms.
func OpenJoystick(path string) (*Joystick, error) { return nil, ErrNotImplemented }

// OpenJoystickDevice is not supported on non-Linux platforms.
func OpenJoystickDevice(info DeviceInfo) (*Joystick, error) { return nil, ErrNotImplemented }

// OpenJoystick is not supported on non-Linux platforms.
func (d DeviceInfo) OpenJoystick() (*Joystick, error) { return nil, ErrNotImplemented }

// Close is not supported on non-Linux platforms.
func (j *Joystick) Close() error { return ErrNotImplemented }

// FD is not supported on non-Linux platforms.
func (j *Joystick) FD() (uintptr, error) { return 0, ErrNotImplemented }

// ReadOnly is not supported on non-Linux platforms.
func (j *Joystick) ReadOnly() bool { return true }

// Version is not supported on non-Linux platforms.
func (j *Joystick) Version() (uint32, error) { return 0, ErrNotImplemented }

// Axes is not supported on non-Linux platforms.
func (j *Joystick) Axes() (uint8, error) { return 0, ErrNotImplemented }

// Buttons is not supported on non-Linux platforms.
func (j *Joystick) Buttons() (uint8, error) { return 0, ErrNotImplemented }

// Name is not supported on non-Linux platforms.
func (j *Joystick) Name() (string, error) { return "", ErrNotImplemented }

// AxisMap is not supported on non-Linux platforms.
func (j *Joystick) AxisMap() ([]uint8, error) { return nil, ErrNotImplemented }

// SetAxisMap is not supported on non-Linux platforms.
func (j *Joystick) SetAxisMap(mapping []uint8) error { return ErrNotImplemented }

// ButtonMap is not supported on non-Linux platforms.
func (j *Joystick) ButtonMap() ([]uint16, error) { return nil, ErrNotImplemented }

// SetButtonMap is not supported on non-Linux platforms.
func (j *Joystick) SetButtonMap(mapping []uint16) error { return ErrNotImplemented }

// Correction is not supported on non-Linux platforms.
func (j *Joystick) Correction() (JoystickCorrection, error) {
	return JoystickCorrection{}, ErrNotImplemented
}

// SetCorrection is not supported on non-Linux platforms.
func (j *Joystick) SetCorrection(corr JoystickCorrection) error { return ErrNotImplemented }

// ReadEvent is not supported on non-Linux platforms.
func (j *Joystick) ReadEvent(timeout time.Duration) (JoystickEvent, error) {
	return JoystickEvent{}, ErrNotImplemented
}
