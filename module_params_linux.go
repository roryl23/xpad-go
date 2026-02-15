//go:build linux

package xpad

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModuleParams represents tunable parameters of the xpad kernel module.
type ModuleParams struct {
	DpadToButtons     bool
	TriggersToButtons bool
	SticksToNull      bool
	AutoPowerOff      bool
}

const moduleParamDir = "/sys/module/xpad/parameters"

// GetModuleParams reads the current xpad module parameters.
func GetModuleParams() (ModuleParams, error) {
	params := ModuleParams{}
	var err error

	if params.DpadToButtons, err = readBoolParam(paramPath("dpad_to_buttons")); err != nil {
		return ModuleParams{}, err
	}
	if params.TriggersToButtons, err = readBoolParam(paramPath("triggers_to_buttons")); err != nil {
		return ModuleParams{}, err
	}
	if params.SticksToNull, err = readBoolParam(paramPath("sticks_to_null")); err != nil {
		return ModuleParams{}, err
	}
	if params.AutoPowerOff, err = readBoolParam(paramPath("auto_poweroff")); err != nil {
		return ModuleParams{}, err
	}

	return params, nil
}

// SetModuleParams writes all xpad module parameters.
//
// This typically requires elevated permissions.
func SetModuleParams(params ModuleParams) error {
	if err := writeBoolParam(paramPath("dpad_to_buttons"), params.DpadToButtons); err != nil {
		return err
	}
	if err := writeBoolParam(paramPath("triggers_to_buttons"), params.TriggersToButtons); err != nil {
		return err
	}
	if err := writeBoolParam(paramPath("sticks_to_null"), params.SticksToNull); err != nil {
		return err
	}
	if err := writeBoolParam(paramPath("auto_poweroff"), params.AutoPowerOff); err != nil {
		return err
	}
	return nil
}

// GetDpadToButtons returns the dpad_to_buttons module parameter.
func GetDpadToButtons() (bool, error) {
	return readBoolParam(paramPath("dpad_to_buttons"))
}

// SetDpadToButtons sets the dpad_to_buttons module parameter.
func SetDpadToButtons(value bool) error {
	return writeBoolParam(paramPath("dpad_to_buttons"), value)
}

// GetTriggersToButtons returns the triggers_to_buttons module parameter.
func GetTriggersToButtons() (bool, error) {
	return readBoolParam(paramPath("triggers_to_buttons"))
}

// SetTriggersToButtons sets the triggers_to_buttons module parameter.
func SetTriggersToButtons(value bool) error {
	return writeBoolParam(paramPath("triggers_to_buttons"), value)
}

// GetSticksToNull returns the sticks_to_null module parameter.
func GetSticksToNull() (bool, error) {
	return readBoolParam(paramPath("sticks_to_null"))
}

// SetSticksToNull sets the sticks_to_null module parameter.
func SetSticksToNull(value bool) error {
	return writeBoolParam(paramPath("sticks_to_null"), value)
}

// GetAutoPowerOff returns the auto_poweroff module parameter.
func GetAutoPowerOff() (bool, error) {
	return readBoolParam(paramPath("auto_poweroff"))
}

// SetAutoPowerOff sets the auto_poweroff module parameter.
func SetAutoPowerOff(value bool) error {
	return writeBoolParam(paramPath("auto_poweroff"), value)
}

func paramPath(name string) string {
	return filepath.Join(moduleParamDir, name)
}

func readBoolParam(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	value := strings.TrimSpace(string(data))
	switch strings.ToLower(value) {
	case "y", "yes", "1", "true", "on":
		return true, nil
	case "n", "no", "0", "false", "off":
		return false, nil
	default:
		return false, fmt.Errorf("xpad: unexpected boolean value %q in %s", value, path)
	}
}

func writeBoolParam(path string, value bool) error {
	text := "0"
	if value {
		text = "1"
	}
	return os.WriteFile(path, []byte(text), 0o644)
}
