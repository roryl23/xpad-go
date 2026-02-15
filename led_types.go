package xpad

// LEDCommand identifies xpad LED patterns.
type LEDCommand uint8

const (
	LEDOff          LEDCommand = 0
	LEDAllBlink     LEDCommand = 1
	LEDPlayer1Flash LEDCommand = 2
	LEDPlayer2Flash LEDCommand = 3
	LEDPlayer3Flash LEDCommand = 4
	LEDPlayer4Flash LEDCommand = 5
	LEDPlayer1      LEDCommand = 6
	LEDPlayer2      LEDCommand = 7
	LEDPlayer3      LEDCommand = 8
	LEDPlayer4      LEDCommand = 9
	LEDRotate1      LEDCommand = 10
	LEDRotate2      LEDCommand = 11
	LEDRotate3      LEDCommand = 12
	LEDRotate4      LEDCommand = 13
	LEDBlinkFast    LEDCommand = 14
	LEDBlinkSlow    LEDCommand = 15
)

// Aliases for xpad LED patterns (commands 10-15) as documented by the driver.
const (
	LEDRotate            LEDCommand = 10
	LEDBlinkPrevious     LEDCommand = 11
	LEDBlinkSlowPrevious LEDCommand = 12
	LEDRotateDual        LEDCommand = 13
	LEDBlinkAllSlow      LEDCommand = 14
	LEDBlinkOnce         LEDCommand = 15
)
