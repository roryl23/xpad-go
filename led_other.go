//go:build !linux

package xpad

// LED represents an xpad LED sysfs device.
type LED struct{}

// OpenLED is not supported on non-Linux platforms.
func OpenLED(path string) (*LED, error) { return nil, ErrNotImplemented }

// OpenLEDDevice is not supported on non-Linux platforms.
func OpenLEDDevice(info DeviceInfo) (*LED, error) { return nil, ErrNotImplemented }

// LED is not supported on non-Linux platforms.
func (d DeviceInfo) LED() (*LED, error) { return nil, ErrNotImplemented }

// SetLED is not supported on non-Linux platforms.
func (d DeviceInfo) SetLED(cmd LEDCommand) error { return ErrNotImplemented }

// SetCommand is not supported on non-Linux platforms.
func (l *LED) SetCommand(cmd LEDCommand) error { return ErrNotImplemented }

// SetBrightness is not supported on non-Linux platforms.
func (l *LED) SetBrightness(value int) error { return ErrNotImplemented }

// Brightness is not supported on non-Linux platforms.
func (l *LED) Brightness() (int, error) { return 0, ErrNotImplemented }
