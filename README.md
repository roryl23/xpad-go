# xpad-go

Go bindings for the Linux kernel `xpad` driver (evdev, joystick, sysfs).

## Support

- Wired Xbox 360 controller only (this is the only device tested/supported today).
- PRs welcome for broader controller support and other features.

Let me reiterate and clarify; This software was written in an afternoon and barely covers a single use case,
that is, button presses on the Xbox 360 controller.
Use at your own risk. If you have a problem feel free to open an issue, and I may get to it in my free time.

Basically, under no circumstances can I be expected to add support for other use cases.
If you need something implemented, be the change you want to see in the world and open a PR!

## Install the xpad driver (Linux)

The `xpad` driver ships with the Linux kernel. In most distros it is already
available as a loadable kernel module.

1. Plug in the Xbox 360 controller (USB or wireless receiver).
2. Load the module if it is not already loaded:

```bash
sudo modprobe xpad
```

3. Verify the devices are present:

```bash
ls -l /dev/input/event*
ls -l /dev/input/js*
ls -l /sys/class/leds/xpad*
```

If you do not see the devices, confirm the kernel module is loaded:

```bash
lsmod | grep xpad
```

Permissions: reading `/dev/input/*` and writing to `/sys/class/leds/*` typically
require elevated privileges or a udev rule that grants your user access.

## Install the Go library

```bash
go get github.com/roryl23/xpad-go
```

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

## Tests

The integration tests require a controller connected on Linux.
If one is not present, the tests that require it will be skipped.

```bash
XPAD_INTERACTIVE=1 go test -run TestIntegrationControllerButtons -v
```

Use `XPAD_EVENT_PATH` or `XPAD_JOYSTICK_PATH` to point tests at a specific
controller if discovery fails.

## Resources

- Linux Kernel driver: https://github.com/paroj/xpad
