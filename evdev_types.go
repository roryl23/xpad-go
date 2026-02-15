package xpad

// Event type constants (EV_*).
const (
	EVSyn      EventKind = 0x00
	EVKey      EventKind = 0x01
	EVRel      EventKind = 0x02
	EVAbs      EventKind = 0x03
	EVMsc      EventKind = 0x04
	EVSw       EventKind = 0x05
	EVLed      EventKind = 0x11
	EVSnd      EventKind = 0x12
	EVRep      EventKind = 0x14
	EVFF       EventKind = 0x15
	EVPwr      EventKind = 0x16
	EVFFStatus EventKind = 0x17
)

// Max code values for common evdev categories.
const (
	EVMax   = 0x1f
	KeyMax  = 0x2ff
	AbsMax  = 0x3f
	AbsCnt  = AbsMax + 1
	FFMax   = 0x7f
	LEDMax  = 0x0f
	BtnMisc = 0x100
)

// Sync event codes (SYN_*).
const (
	SynReport   = 0
	SynConfig   = 1
	SynMTReport = 2
	SynDropped  = 3
)

// Absolute axis codes (ABS_*).
const (
	ABSX       = 0x00
	ABSY       = 0x01
	ABSZ       = 0x02
	ABSRX      = 0x03
	ABSRY      = 0x04
	ABSRZ      = 0x05
	ABSHat0X   = 0x10
	ABSHat0Y   = 0x11
	ABSProfile = 0x21
)

// Button codes for common Xbox controller inputs (BTN_*).
const (
	BTNA      = 0x130
	BTNB      = 0x131
	BTNX      = 0x133
	BTNY      = 0x134
	BTNTL     = 0x136
	BTNTR     = 0x137
	BTNTL2    = 0x138
	BTNTR2    = 0x139
	BTNSelect = 0x13a
	BTNStart  = 0x13b
	BTNMode   = 0x13c
	BTNThumbL = 0x13d
	BTNThumbR = 0x13e
)

// D-pad button codes (BTN_DPAD_*), used when dpad_to_buttons is enabled.
const (
	BTNDPadUp    = 0x220
	BTNDPadDown  = 0x221
	BTNDPadLeft  = 0x222
	BTNDPadRight = 0x223
)

// Trigger-happy button codes (BTN_TRIGGER_HAPPY*), used by xpad for dpad/paddles.
const (
	BTNTriggerHappy1 = 0x2c0
	BTNTriggerHappy2 = 0x2c1
	BTNTriggerHappy3 = 0x2c2
	BTNTriggerHappy4 = 0x2c3
	BTNTriggerHappy5 = 0x2c4
	BTNTriggerHappy6 = 0x2c5
	BTNTriggerHappy7 = 0x2c6
	BTNTriggerHappy8 = 0x2c7
)

// Key codes used by xpad for select/share mappings.
const (
	KeyRecord = 0x0a7
)

// Force feedback effect types (FF_*).
const (
	FFRumble = 0x50
)

// Force feedback device properties.
const (
	FFGain       = 0x60
	FFAutocenter = 0x61
)

// InputID mirrors struct input_id.
type InputID struct {
	BusType uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

// AbsInfo mirrors struct input_absinfo.
type AbsInfo struct {
	Value      int32
	Minimum    int32
	Maximum    int32
	Fuzz       int32
	Flat       int32
	Resolution int32
}
