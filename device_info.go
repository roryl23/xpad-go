package xpad

import "strings"

// DeviceInfo describes a discovered input device.
type DeviceInfo struct {
	// Path is the event device path, typically /dev/input/eventX.
	Path string
	// SysfsPath is the sysfs node for the event device, typically /sys/class/input/eventX.
	SysfsPath string

	// DevicePath is the resolved sysfs path to the backing device.
	DevicePath string
	// JoystickPath is the matching /dev/input/jsX device if present.
	JoystickPath string
	// LEDPath is the sysfs directory for the xpad LED device if present.
	LEDPath string
	// LEDBrightnessPath is the sysfs brightness file for the LED device if present.
	LEDBrightnessPath string

	Name      string
	Phys      string
	Uniq      string
	Driver    string
	BusType   uint16
	VendorID  uint16
	ProductID uint16
	VersionID uint16
}

// IsXpad reports whether the device appears to be handled by the xpad driver.
func (d DeviceInfo) IsXpad() bool {
	driver := strings.ToLower(d.Driver)
	if strings.Contains(driver, "xpad") {
		return true
	}

	name := strings.ToLower(d.Name)
	return strings.Contains(name, "xpad") || strings.Contains(name, "xbox") || strings.Contains(name, "x-box")
}
