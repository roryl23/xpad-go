//go:build linux

package xpad

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ListDevices scans /dev/input for event devices and enriches them via sysfs.
func ListDevices() ([]DeviceInfo, error) {
	jsMap, err := mapSysfsDevices("/sys/class/input/js*", "/dev/input")
	if err != nil {
		return nil, err
	}
	ledMap, err := mapSysfsDevices("/sys/class/leds/xpad*", "")
	if err != nil {
		return nil, err
	}

	paths, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return nil, err
	}

	infos := make([]DeviceInfo, 0, len(paths))
	for _, path := range paths {
		base := filepath.Base(path)
		sysfs := filepath.Join("/sys/class/input", base)

		info := DeviceInfo{Path: path, SysfsPath: sysfs}
		devPath := filepath.Join(sysfs, "device")
		if resolved, err := filepath.EvalSymlinks(devPath); err == nil {
			info.DevicePath = resolved
			if jsPath, ok := jsMap[resolved]; ok {
				info.JoystickPath = jsPath
			}
			if ledPath, ok := ledMap[resolved]; ok {
				info.LEDPath = ledPath
				info.LEDBrightnessPath = filepath.Join(ledPath, "brightness")
			}
		}

		info.Name = readTrimmedFile(filepath.Join(devPath, "name"))
		info.Phys = readTrimmedFile(filepath.Join(devPath, "phys"))
		info.Uniq = readTrimmedFile(filepath.Join(devPath, "uniq"))
		info.Driver = readLinkBase(filepath.Join(devPath, "driver"))

		if v, ok := readHexUint16(filepath.Join(devPath, "id", "bustype")); ok {
			info.BusType = v
		}
		if v, ok := readHexUint16(filepath.Join(devPath, "id", "vendor")); ok {
			info.VendorID = v
		}
		if v, ok := readHexUint16(filepath.Join(devPath, "id", "product")); ok {
			info.ProductID = v
		}
		if v, ok := readHexUint16(filepath.Join(devPath, "id", "version")); ok {
			info.VersionID = v
		}

		infos = append(infos, info)
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Path < infos[j].Path
	})

	return infos, nil
}

// FindXpadDevices returns only devices that look like xpad-backed controllers.
func FindXpadDevices() ([]DeviceInfo, error) {
	infos, err := ListDevices()
	if err != nil {
		return nil, err
	}
	filtered := make([]DeviceInfo, 0, len(infos))
	for _, info := range infos {
		if info.IsXpad() {
			filtered = append(filtered, info)
		}
	}
	return filtered, nil
}

// OpenFirstXpad opens the first xpad-backed device discovered.
func OpenFirstXpad() (*Device, error) {
	infos, err := FindXpadDevices()
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, ErrNotFound
	}
	return Open(infos[0].Path)
}

func readTrimmedFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readHexUint16(path string) (uint16, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return 0, false
	}
	value, err := strconv.ParseUint(text, 16, 16)
	if err != nil {
		value, err = strconv.ParseUint(text, 0, 16)
		if err != nil {
			return 0, false
		}
	}
	return uint16(value), true
}

func readLinkBase(path string) string {
	link, err := os.Readlink(path)
	if err != nil {
		return ""
	}
	return filepath.Base(link)
}

func mapSysfsDevices(globPattern, devPrefix string) (map[string]string, error) {
	paths, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, err
	}
	mapping := make(map[string]string, len(paths))
	for _, sysfs := range paths {
		deviceLink := filepath.Join(sysfs, "device")
		resolved, err := filepath.EvalSymlinks(deviceLink)
		if err != nil {
			continue
		}
		base := filepath.Base(sysfs)
		if devPrefix != "" {
			mapping[resolved] = filepath.Join(devPrefix, base)
		} else {
			mapping[resolved] = sysfs
		}
	}
	return mapping, nil
}
