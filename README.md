# xpad-go

Go bindings for the Linux kernel `xpad` driver (evdev, joystick, sysfs).

Linux only.

Features:
- Device discovery and open/close helpers for `/dev/input/event*`
- Mapping to `/dev/input/js*` and `/sys/class/leds/xpad*`
- evdev metadata, capability queries, and event I/O
- Force-feedback rumble upload/play/erase
- LED control via sysfs brightness
- Module parameter read/write helpers

## Quick start

```go
devices, err := xpad.FindXpadDevices()
if err != nil {
	// handle error
}
if len(devices) == 0 {
	// no controllers found
}

dev, err := xpad.OpenDevice(devices[0])
if err != nil {
	// handle error
}
defer dev.Close()

event, err := dev.ReadEvent(-1)
if err != nil {
	// handle error
}
_ = event
```

## LED control

```go
info := devices[0]
if err := info.SetLED(xpad.LEDPlayer1); err != nil {
	// handle error
}
```

LED command values:
- 0: off
- 1: all blink, then previous setting
- 2-5: 1-4 blink, then on
- 6-9: 1-4 on
- 10: rotate
- 11: blink based on previous setting
- 12: slow blink based on previous setting
- 13: rotate with two lights
- 14: persistent slow all blink
- 15: blink once, then previous setting

## Rumble

```go
id, err := dev.UploadRumble(xpad.NewRumbleEffect(0xffff, 0x7fff, 500*time.Millisecond))
if err != nil {
	// handle error
}
if err := dev.PlayEffect(id, 1); err != nil {
	// handle error
}
```

## Joystick API

```go
js, err := info.OpenJoystick()
if err != nil {
	// handle error
}
defer js.Close()

evt, err := js.ReadEvent(-1)
if err != nil {
	// handle error
}
_ = evt
```
