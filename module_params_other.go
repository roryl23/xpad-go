//go:build !linux

package xpad

// ModuleParams represents tunable parameters of the xpad kernel module.
type ModuleParams struct {
	DpadToButtons     bool
	TriggersToButtons bool
	SticksToNull      bool
	AutoPowerOff      bool
}

// GetModuleParams is not supported on non-Linux platforms.
func GetModuleParams() (ModuleParams, error) {
	return ModuleParams{}, ErrNotImplemented
}

// SetModuleParams is not supported on non-Linux platforms.
func SetModuleParams(params ModuleParams) error {
	return ErrNotImplemented
}

// GetDpadToButtons is not supported on non-Linux platforms.
func GetDpadToButtons() (bool, error) {
	return false, ErrNotImplemented
}

// SetDpadToButtons is not supported on non-Linux platforms.
func SetDpadToButtons(value bool) error {
	return ErrNotImplemented
}

// GetTriggersToButtons is not supported on non-Linux platforms.
func GetTriggersToButtons() (bool, error) {
	return false, ErrNotImplemented
}

// SetTriggersToButtons is not supported on non-Linux platforms.
func SetTriggersToButtons(value bool) error {
	return ErrNotImplemented
}

// GetSticksToNull is not supported on non-Linux platforms.
func GetSticksToNull() (bool, error) {
	return false, ErrNotImplemented
}

// SetSticksToNull is not supported on non-Linux platforms.
func SetSticksToNull(value bool) error {
	return ErrNotImplemented
}

// GetAutoPowerOff is not supported on non-Linux platforms.
func GetAutoPowerOff() (bool, error) {
	return false, ErrNotImplemented
}

// SetAutoPowerOff is not supported on non-Linux platforms.
func SetAutoPowerOff(value bool) error {
	return ErrNotImplemented
}
