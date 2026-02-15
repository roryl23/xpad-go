//go:build linux

package xpad

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func TestIntegrationEvdevBasics(t *testing.T) {
	info, ok := findTestDeviceInfo(t)
	if !ok {
		return
	}

	dev, err := OpenDevice(info)
	if err != nil {
		t.Skipf("open device failed: %v", err)
	}
	defer dev.Close()

	if _, err := dev.Name(); err != nil {
		t.Fatalf("Name() error: %v", err)
	}
	if _, err := dev.Phys(); err != nil {
		if isOptionalDeviceInfoError(err) {
			t.Logf("Phys() not supported: %v", err)
		} else {
			t.Fatalf("Phys() error: %v", err)
		}
	}
	if _, err := dev.Uniq(); err != nil {
		if isOptionalDeviceInfoError(err) {
			t.Logf("Uniq() not supported: %v", err)
		} else {
			t.Fatalf("Uniq() error: %v", err)
		}
	}
	if _, err := dev.ID(); err != nil {
		if isOptionalEvdevError(err) {
			t.Skipf("ID() not supported: %v", err)
		}
		t.Fatalf("ID() error: %v", err)
	}

	types, err := dev.EventTypes()
	if err != nil {
		if isOptionalEvdevError(err) {
			t.Skipf("EventTypes() not supported: %v", err)
		}
		t.Fatalf("EventTypes() error: %v", err)
	}
	if len(types) == 0 {
		t.Fatalf("EventTypes() returned empty bitset")
	}

	hasAbs, err := dev.HasEventType(EVAbs)
	if err != nil {
		if isOptionalEvdevError(err) {
			t.Skipf("HasEventType(EVAbs) not supported: %v", err)
		}
		t.Fatalf("HasEventType(EVAbs) error: %v", err)
	}
	if hasAbs {
		hasX, err := dev.HasEventCode(EVAbs, ABSX)
		if err != nil {
			if isOptionalEvdevError(err) {
				t.Skipf("HasEventCode(EVAbs, ABSX) not supported: %v", err)
			}
			t.Fatalf("HasEventCode(EVAbs, ABSX) error: %v", err)
		}
		if hasX {
			if _, err := dev.AbsInfo(ABSX); err != nil {
				if isOptionalEvdevError(err) {
					t.Skipf("AbsInfo(ABSX) not supported: %v", err)
				}
				t.Fatalf("AbsInfo(ABSX) error: %v", err)
			}
		}
	}

	if _, err := dev.ReadEvent(5 * time.Millisecond); err != nil && !errors.Is(err, ErrTimeout) {
		t.Fatalf("ReadEvent() error: %v", err)
	}
}

func TestIntegrationJoystickBasics(t *testing.T) {
	info, ok := findTestDeviceInfo(t)
	if !ok {
		return
	}
	if info.JoystickPath == "" {
		t.Skip("no joystick device discovered")
	}

	js, err := OpenJoystickDevice(info)
	if err != nil {
		t.Skipf("open joystick failed: %v", err)
	}
	defer js.Close()

	if _, err := js.Version(); err != nil {
		if isOptionalJoystickError(err) {
			t.Skipf("joystick ioctls not supported: %v", err)
		}
		t.Fatalf("Version() error: %v", err)
	}
	if _, err := js.Name(); err != nil {
		t.Fatalf("Name() error: %v", err)
	}
	if _, err := js.Axes(); err != nil {
		if isOptionalJoystickError(err) {
			t.Skipf("Axes() not supported: %v", err)
		}
		t.Fatalf("Axes() error: %v", err)
	}
	if _, err := js.Buttons(); err != nil {
		if isOptionalJoystickError(err) {
			t.Skipf("Buttons() not supported: %v", err)
		}
		t.Fatalf("Buttons() error: %v", err)
	}
	if _, err := js.AxisMap(); err != nil {
		if isOptionalJoystickError(err) {
			t.Skipf("AxisMap() not supported: %v", err)
		}
		t.Fatalf("AxisMap() error: %v", err)
	}
	if _, err := js.ButtonMap(); err != nil {
		if isOptionalJoystickError(err) {
			t.Skipf("ButtonMap() not supported: %v", err)
		}
		t.Fatalf("ButtonMap() error: %v", err)
	}

	if _, err := js.ReadEvent(5 * time.Millisecond); err != nil && !errors.Is(err, ErrTimeout) {
		t.Fatalf("ReadEvent() error: %v", err)
	}
}

func TestIntegrationLEDBasics(t *testing.T) {
	info, ok := findTestDeviceInfo(t)
	if !ok {
		return
	}
	if info.LEDPath == "" && info.LEDBrightnessPath == "" {
		t.Skip("no LED sysfs device discovered")
	}

	led, err := OpenLEDDevice(info)
	if err != nil {
		t.Skipf("open LED failed: %v", err)
	}

	value, err := led.Brightness()
	if err != nil {
		if os.IsPermission(err) || os.IsNotExist(err) {
			t.Skipf("brightness read not permitted: %v", err)
		}
		t.Fatalf("Brightness() error: %v", err)
	}
	if value < 0 || value > 15 {
		t.Fatalf("Brightness() = %d, expected range 0-15", value)
	}

	if os.Getenv("XPAD_TEST_WRITE_LED") != "" {
		if err := led.SetBrightness(value); err != nil {
			t.Fatalf("SetBrightness() error: %v", err)
		}
	}
}

func TestIntegrationModuleParams(t *testing.T) {
	_, ok := findTestDeviceInfo(t)
	if !ok {
		return
	}
	if _, err := GetModuleParams(); err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			t.Skipf("module params unavailable: %v", err)
		}
		t.Fatalf("GetModuleParams() error: %v", err)
	}
}

func TestIntegrationControllerButtons(t *testing.T) {
	if !requireInteractive(t) {
		return
	}

	profile := loadControllerProfile(t)
	timeout := interactiveTimeout()

	t.Run(profile.Name+"/evdev-buttons", func(t *testing.T) {
		info, ok := findTestDeviceInfo(t)
		if !ok {
			return
		}
		dev, err := OpenDevice(info)
		if err != nil {
			t.Skipf("open device failed: %v", err)
		}
		defer dev.Close()

		hasKey, err := dev.HasEventType(EVKey)
		if err != nil {
			if isOptionalEvdevError(err) {
				t.Skipf("HasEventType(EVKey) not supported: %v", err)
			}
			t.Fatalf("HasEventType(EVKey) error: %v", err)
		}
		if !hasKey {
			t.Skip("device does not report EV_KEY support")
		}

		t.Logf("Controller profile %q: press each button when prompted (timeout %s each).", profile.Name, timeout)
		for _, step := range profile.EvdevSteps {
			t.Logf("Press %s", step.Name)
			if err := waitForEvdevStep(dev, step, timeout); err != nil {
				if errors.Is(err, ErrTimeout) {
					t.Fatalf("timed out waiting for %s", step.Name)
				}
				t.Fatalf("waiting for %s failed: %v", step.Name, err)
			}
		}
	})

	t.Run(profile.Name+"/joystick-buttons", func(t *testing.T) {
		info, ok := findTestDeviceInfo(t)
		if !ok {
			return
		}
		if info.JoystickPath == "" {
			t.Skip("no joystick device discovered")
		}

		js, err := OpenJoystickDevice(info)
		if err != nil {
			t.Skipf("open joystick failed: %v", err)
		}
		defer js.Close()

		if _, err := js.Version(); err != nil {
			if isOptionalJoystickError(err) {
				t.Skipf("joystick ioctls not supported: %v", err)
			}
			t.Fatalf("Version() error: %v", err)
		}

		buttonMap, err := js.ButtonMap()
		if err != nil {
			if isOptionalJoystickError(err) {
				t.Skipf("ButtonMap() not supported: %v", err)
			}
			t.Fatalf("ButtonMap() error: %v", err)
		}
		joyButtons := buildJoystickButtonMap(buttonMap)
		if len(profile.JoystickButtons) == 0 {
			t.Skip("no joystick button expectations defined for this profile")
		}

		for _, step := range profile.JoystickButtons {
			button, ok := joystickButtonForKeys(joyButtons, step.KeyCodes)
			if !ok {
				t.Skipf("joystick mapping missing for %s", step.Name)
			}
			t.Logf("Press %s", step.Name)
			if err := waitForJoystickButton(js, button, timeout); err != nil {
				if errors.Is(err, ErrTimeout) {
					t.Fatalf("timed out waiting for %s", step.Name)
				}
				t.Fatalf("waiting for %s failed: %v", step.Name, err)
			}
		}
	})
}

func findTestDeviceInfo(t *testing.T) (DeviceInfo, bool) {
	t.Helper()

	if path := os.Getenv("XPAD_EVENT_PATH"); path != "" {
		info := DeviceInfo{Path: path}
		if jsPath := os.Getenv("XPAD_JOYSTICK_PATH"); jsPath != "" {
			info.JoystickPath = jsPath
		}
		if enriched, ok := enrichDeviceInfo(path); ok {
			if info.JoystickPath != "" {
				enriched.JoystickPath = info.JoystickPath
			}
			info = enriched
		}
		return info, true
	}

	infos, err := FindXpadDevices()
	if err != nil {
		t.Skipf("FindXpadDevices() failed: %v", err)
		return DeviceInfo{}, false
	}
	if len(infos) == 0 {
		if info, ok := findDeviceInfoByJoystick(t); ok {
			return info, true
		}
		logAvailableDevices(t)
		t.Skip("no xpad devices found")
		return DeviceInfo{}, false
	}
	return infos[0], true
}

func enrichDeviceInfo(path string) (DeviceInfo, bool) {
	infos, err := ListDevices()
	if err != nil {
		return DeviceInfo{}, false
	}
	for _, info := range infos {
		if info.Path == path {
			return info, true
		}
	}
	return DeviceInfo{}, false
}

func findDeviceInfoByJoystick(t *testing.T) (DeviceInfo, bool) {
	t.Helper()

	if jsPath := os.Getenv("XPAD_JOYSTICK_PATH"); jsPath != "" {
		info := DeviceInfo{JoystickPath: jsPath}
		if enriched, ok := enrichDeviceInfoByJoystick(jsPath); ok {
			info = enriched
		}
		return info, true
	}

	jsPaths, err := filepath.Glob("/dev/input/js*")
	if err != nil {
		t.Logf("failed to glob /dev/input/js*: %v", err)
		return DeviceInfo{}, false
	}
	if len(jsPaths) == 0 {
		t.Log("no joystick devices found under /dev/input/js*")
		return DeviceInfo{}, false
	}

	infos, err := ListDevices()
	if err != nil {
		t.Logf("ListDevices() failed while mapping joystick devices: %v", err)
		return DeviceInfo{}, false
	}

	for _, jsPath := range jsPaths {
		for _, info := range infos {
			if info.JoystickPath == jsPath {
				t.Logf("using joystick device %s mapped to %s", jsPath, info.Path)
				return info, true
			}
		}
	}

	for _, jsPath := range jsPaths {
		base := filepath.Base(jsPath)
		sysfs := filepath.Join("/sys/class/input", base)
		deviceLink := filepath.Join(sysfs, "device")
		resolved, err := filepath.EvalSymlinks(deviceLink)
		if err != nil {
			continue
		}
		for _, info := range infos {
			if info.DevicePath == resolved && info.Path != "" {
				info.JoystickPath = jsPath
				t.Logf("using joystick device %s mapped via sysfs to %s", jsPath, info.Path)
				return info, true
			}
		}
	}

	t.Log("found joystick devices but could not map them to event devices")
	return DeviceInfo{}, false
}

func enrichDeviceInfoByJoystick(jsPath string) (DeviceInfo, bool) {
	infos, err := ListDevices()
	if err != nil {
		return DeviceInfo{}, false
	}
	for _, info := range infos {
		if info.JoystickPath == jsPath {
			return info, true
		}
	}
	return DeviceInfo{}, false
}

func logAvailableDevices(t *testing.T) {
	t.Helper()

	infos, err := ListDevices()
	if err != nil {
		t.Logf("ListDevices() failed while listing available devices: %v", err)
		return
	}
	if len(infos) == 0 {
		t.Log("ListDevices() returned no devices")
		return
	}
	t.Log("Available input devices:")
	for _, info := range infos {
		name := info.Name
		if name == "" {
			name = "(unknown)"
		}
		driver := info.Driver
		if driver == "" {
			driver = "(unknown)"
		}
		t.Logf("  %s name=%q driver=%q", info.Path, name, driver)
	}
}

func requireInteractive(t *testing.T) bool {
	if os.Getenv("XPAD_INTERACTIVE") == "" {
		t.Skip("set XPAD_INTERACTIVE=1 to enable interactive controller tests")
		return false
	}
	return true
}

type controllerProfile struct {
	Name            string
	EvdevSteps      []evdevStep
	JoystickButtons []buttonStep
}

type evdevStep struct {
	Name     string
	KeyCodes []uint16
	AbsCode  uint16
	AbsMatch func(int32) bool
}

type buttonStep struct {
	Name     string
	KeyCodes []uint16
}

var controllerProfiles = map[string]controllerProfile{
	"xbox360": xbox360Profile(),
}

func loadControllerProfile(t *testing.T) controllerProfile {
	t.Helper()

	name := os.Getenv("XPAD_CONTROLLER")
	if name == "" {
		name = "xbox360"
	}
	profile, ok := controllerProfiles[name]
	if !ok {
		t.Skipf("unknown controller profile %q (set XPAD_CONTROLLER)", name)
	}
	return profile
}

func xbox360Profile() controllerProfile {
	return controllerProfile{
		Name: "xbox360",
		EvdevSteps: []evdevStep{
			{Name: "A", KeyCodes: []uint16{BTNA}},
			{Name: "B", KeyCodes: []uint16{BTNB}},
			{Name: "X", KeyCodes: []uint16{BTNX}},
			{Name: "Y", KeyCodes: []uint16{BTNY}},
			{Name: "LB", KeyCodes: []uint16{BTNTL}},
			{Name: "RB", KeyCodes: []uint16{BTNTR}},
			{Name: "Back", KeyCodes: []uint16{BTNSelect}},
			{Name: "Start", KeyCodes: []uint16{BTNStart}},
			{Name: "Guide", KeyCodes: []uint16{BTNMode}},
			{Name: "Left Stick Click", KeyCodes: []uint16{BTNThumbL}},
			{Name: "Right Stick Click", KeyCodes: []uint16{BTNThumbR}},
			{
				Name:     "D-pad Up",
				KeyCodes: []uint16{BTNDPadUp},
				AbsCode:  ABSHat0Y,
				AbsMatch: func(v int32) bool { return v < 0 },
			},
			{
				Name:     "D-pad Down",
				KeyCodes: []uint16{BTNDPadDown},
				AbsCode:  ABSHat0Y,
				AbsMatch: func(v int32) bool { return v > 0 },
			},
			{
				Name:     "D-pad Left",
				KeyCodes: []uint16{BTNDPadLeft},
				AbsCode:  ABSHat0X,
				AbsMatch: func(v int32) bool { return v < 0 },
			},
			{
				Name:     "D-pad Right",
				KeyCodes: []uint16{BTNDPadRight},
				AbsCode:  ABSHat0X,
				AbsMatch: func(v int32) bool { return v > 0 },
			},
			{
				Name:     "Left Trigger",
				KeyCodes: []uint16{BTNTL2},
				AbsCode:  ABSZ,
				AbsMatch: func(v int32) bool { return v > 10 },
			},
			{
				Name:     "Right Trigger",
				KeyCodes: []uint16{BTNTR2},
				AbsCode:  ABSRZ,
				AbsMatch: func(v int32) bool { return v > 10 },
			},
		},
		JoystickButtons: []buttonStep{
			{Name: "A", KeyCodes: []uint16{BTNA}},
			{Name: "B", KeyCodes: []uint16{BTNB}},
			{Name: "X", KeyCodes: []uint16{BTNX}},
			{Name: "Y", KeyCodes: []uint16{BTNY}},
			{Name: "LB", KeyCodes: []uint16{BTNTL}},
			{Name: "RB", KeyCodes: []uint16{BTNTR}},
			{Name: "Back", KeyCodes: []uint16{BTNSelect}},
			{Name: "Start", KeyCodes: []uint16{BTNStart}},
			{Name: "Guide", KeyCodes: []uint16{BTNMode}},
			{Name: "Left Stick Click", KeyCodes: []uint16{BTNThumbL}},
			{Name: "Right Stick Click", KeyCodes: []uint16{BTNThumbR}},
		},
	}
}

func waitForEvdevStep(dev *Device, step evdevStep, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ev, err := dev.ReadEvent(500 * time.Millisecond)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				continue
			}
			return err
		}
		if matchEvdevStep(step, ev) {
			return nil
		}
	}
	return ErrTimeout
}

func matchEvdevStep(step evdevStep, ev Event) bool {
	switch ev.Kind {
	case EVKey:
		if ev.Value == 0 {
			return false
		}
		for _, code := range step.KeyCodes {
			if ev.Code == code {
				return true
			}
		}
	case EVAbs:
		if step.AbsCode == 0 || ev.Code != step.AbsCode {
			return false
		}
		if step.AbsMatch == nil {
			return ev.Value != 0
		}
		return step.AbsMatch(ev.Value)
	}
	return false
}

func waitForJoystickButton(js *Joystick, button uint8, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ev, err := js.ReadEvent(500 * time.Millisecond)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				continue
			}
			return err
		}
		kind := ev.Type &^ JoyEventInit
		if kind == JoyEventButton && ev.Number == button && ev.Value != 0 {
			return nil
		}
	}
	return ErrTimeout
}

func buildJoystickButtonMap(mapping []uint16) map[uint16]uint8 {
	buttons := make(map[uint16]uint8)
	for idx, code := range mapping {
		if code == 0 {
			continue
		}
		if _, exists := buttons[code]; !exists {
			buttons[code] = uint8(idx)
		}
	}
	return buttons
}

func joystickButtonForKeys(buttons map[uint16]uint8, codes []uint16) (uint8, bool) {
	for _, code := range codes {
		if btn, ok := buttons[code]; ok {
			return btn, true
		}
	}
	return 0, false
}

func interactiveTimeout() time.Duration {
	if value := os.Getenv("XPAD_INTERACTIVE_TIMEOUT"); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return 10 * time.Second
}

func isOptionalDeviceInfoError(err error) bool {
	return errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ENOTTY) || os.IsNotExist(err)
}

func isOptionalJoystickError(err error) bool {
	return errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.ENOTTY) || errors.Is(err, syscall.ENODEV) || os.IsNotExist(err)
}

func isOptionalEvdevError(err error) bool {
	return errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.ENOTTY) || errors.Is(err, syscall.ENODEV) || os.IsNotExist(err)
}
