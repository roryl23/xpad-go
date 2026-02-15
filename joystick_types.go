package xpad

// JoystickEventType identifies the kind of joystick event.
type JoystickEventType uint8

const (
	JoyEventButton JoystickEventType = 0x01
	JoyEventAxis   JoystickEventType = 0x02
	JoyEventInit   JoystickEventType = 0x80
)

// JoystickEvent mirrors struct js_event from /dev/input/js*.
type JoystickEvent struct {
	Time   uint32
	Value  int16
	Type   JoystickEventType
	Number uint8
}

// JoystickCorrectionType identifies the correction type.
type JoystickCorrectionType uint16

const (
	JoyCorrNone   JoystickCorrectionType = 0x00
	JoyCorrBroken JoystickCorrectionType = 0x01
)

// JoystickCorrection mirrors struct js_corr.
type JoystickCorrection struct {
	Coef [8]int32
	Prec int16
	Type JoystickCorrectionType
}
