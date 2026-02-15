//go:build !linux

package xpad

// ListDevices is not supported on non-Linux platforms.
func ListDevices() ([]DeviceInfo, error) {
	return nil, ErrNotImplemented
}

// FindXpadDevices is not supported on non-Linux platforms.
func FindXpadDevices() ([]DeviceInfo, error) {
	return nil, ErrNotImplemented
}

// OpenFirstXpad is not supported on non-Linux platforms.
func OpenFirstXpad() (*Device, error) {
	return nil, ErrNotImplemented
}
