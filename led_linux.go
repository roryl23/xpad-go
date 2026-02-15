//go:build linux

package xpad

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// LED represents an xpad LED sysfs device.
type LED struct {
	Path           string
	BrightnessPath string
}

// OpenLED opens an LED sysfs device by path or brightness file.
func OpenLED(path string) (*LED, error) {
	if path == "" {
		return nil, ErrNotFound
	}
	brightness := path
	if filepath.Base(path) != "brightness" {
		brightness = filepath.Join(path, "brightness")
	}
	return &LED{Path: filepath.Dir(brightness), BrightnessPath: brightness}, nil
}

// OpenLEDDevice opens the LED for a discovered device.
func OpenLEDDevice(info DeviceInfo) (*LED, error) {
	if info.LEDBrightnessPath == "" && info.LEDPath == "" {
		return nil, ErrNotFound
	}
	if info.LEDBrightnessPath != "" {
		return OpenLED(info.LEDBrightnessPath)
	}
	return OpenLED(info.LEDPath)
}

// LED returns the LED for a discovered device.
func (d DeviceInfo) LED() (*LED, error) {
	return OpenLEDDevice(d)
}

// SetLED sets the LED command for a discovered device.
func (d DeviceInfo) SetLED(cmd LEDCommand) error {
	led, err := d.LED()
	if err != nil {
		return err
	}
	return led.SetCommand(cmd)
}

// SetCommand writes a LED command to sysfs.
func (l *LED) SetCommand(cmd LEDCommand) error {
	return l.SetBrightness(int(cmd))
}

// SetBrightness writes a raw brightness value (0-15).
func (l *LED) SetBrightness(value int) error {
	if l == nil || l.BrightnessPath == "" {
		return ErrNotFound
	}
	if value < 0 || value > 15 {
		return fmt.Errorf("xpad: LED brightness out of range: %d", value)
	}
	return os.WriteFile(l.BrightnessPath, []byte(strconv.Itoa(value)), 0o644)
}

// Brightness reads the current LED brightness value.
func (l *LED) Brightness() (int, error) {
	if l == nil || l.BrightnessPath == "" {
		return 0, ErrNotFound
	}
	data, err := os.ReadFile(l.BrightnessPath)
	if err != nil {
		return 0, err
	}
	text := strings.TrimSpace(string(data))
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return value, nil
}
