//go:build linux

package xpad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"syscall"
	"time"
	"unsafe"

	"github.com/roryl23/xpad-go/internal/ioctl"
)

const evdevIOCBase = 0x45 // 'E'

func evioCGNAME(length uint) uint {
	return ioctl.IOC(ioctl.DirRead, evdevIOCBase, 0x06, length)
}

func evioCGPHYS(length uint) uint {
	return ioctl.IOC(ioctl.DirRead, evdevIOCBase, 0x07, length)
}

func evioCGUNIQ(length uint) uint {
	return ioctl.IOC(ioctl.DirRead, evdevIOCBase, 0x08, length)
}

func evioCGID() uint {
	return ioctl.IOR(evdevIOCBase, 0x02, ioctl.Size(InputID{}))
}

func evioCGBIT(ev EventKind, length uint) uint {
	return ioctl.IOC(ioctl.DirRead, evdevIOCBase, uint(0x20)+uint(ev), length)
}

func evioCGABS(code uint16) uint {
	return ioctl.IOR(evdevIOCBase, uint(0x40)+uint(code), ioctl.Size(AbsInfo{}))
}

func evioCSFF() uint {
	return ioctl.IOW(evdevIOCBase, 0x80, ioctl.Size(ffEffect{}))
}

func evioCRMFF() uint {
	return ioctl.IOW(evdevIOCBase, 0x81, ioctl.Size(int32(0)))
}

func evioCGEFFECTS() uint {
	return ioctl.IOR(evdevIOCBase, 0x84, ioctl.Size(int32(0)))
}

func evioCGRAB() uint {
	return ioctl.IOW(evdevIOCBase, 0x90, ioctl.Size(int32(0)))
}

// Name returns the evdev device name.
func (d *Device) Name() (string, error) {
	return getStringIoctl(d, evioCGNAME)
}

// Phys returns the evdev physical path string.
func (d *Device) Phys() (string, error) {
	return getStringIoctl(d, evioCGPHYS)
}

// Uniq returns the evdev unique identifier string.
func (d *Device) Uniq() (string, error) {
	return getStringIoctl(d, evioCGUNIQ)
}

// ID returns the evdev input_id data.
func (d *Device) ID() (InputID, error) {
	if d == nil || d.file == nil {
		return InputID{}, ErrClosed
	}
	fd, _ := d.FD()
	var id InputID
	if err := ioctl.CallPtr(fd, evioCGID(), unsafe.Pointer(&id)); err != nil {
		return InputID{}, err
	}
	return id, nil
}

// AbsInfo returns absolute axis metadata for the provided code.
func (d *Device) AbsInfo(code uint16) (AbsInfo, error) {
	if d == nil || d.file == nil {
		return AbsInfo{}, ErrClosed
	}
	fd, _ := d.FD()
	var info AbsInfo
	if err := ioctl.CallPtr(fd, evioCGABS(code), unsafe.Pointer(&info)); err != nil {
		return AbsInfo{}, err
	}
	return info, nil
}

// EventTypes returns a bitset of supported event types.
func (d *Device) EventTypes() ([]byte, error) {
	return d.eventBitset(0, EVMax)
}

// HasEventType reports whether the device supports the provided event type.
func (d *Device) HasEventType(ev EventKind) (bool, error) {
	bits, err := d.EventTypes()
	if err != nil {
		return false, err
	}
	return bitsetHas(bits, uint16(ev)), nil
}

// HasEventCode reports whether the device supports the provided event code.
func (d *Device) HasEventCode(ev EventKind, code uint16) (bool, error) {
	var max uint16
	switch ev {
	case EVKey:
		max = KeyMax
	case EVAbs:
		max = AbsMax
	case EVFF:
		max = FFMax
	case EVLed:
		max = LEDMax
	default:
		return false, fmt.Errorf("xpad: unsupported event type %d", ev)
	}
	bits, err := d.eventBitset(ev, max)
	if err != nil {
		return false, err
	}
	return bitsetHas(bits, code), nil
}

// EffectCount returns the number of force-feedback effects supported.
func (d *Device) EffectCount() (int, error) {
	if d == nil || d.file == nil {
		return 0, ErrClosed
	}
	fd, _ := d.FD()
	var count int32
	if err := ioctl.CallPtr(fd, evioCGEFFECTS(), unsafe.Pointer(&count)); err != nil {
		return 0, err
	}
	return int(count), nil
}

// Grab enables or disables exclusive access to the device.
func (d *Device) Grab(grab bool) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	fd, _ := d.FD()
	var value int32
	if grab {
		value = 1
	}
	return ioctl.CallPtr(fd, evioCGRAB(), unsafe.Pointer(&value))
}

func (d *Device) eventBitset(ev EventKind, max uint16) ([]byte, error) {
	if d == nil || d.file == nil {
		return nil, ErrClosed
	}
	length := bitsetBytes(max)
	buf := make([]byte, length)
	fd, _ := d.FD()
	if err := ioctl.CallPtr(fd, evioCGBIT(ev, uint(length)), unsafe.Pointer(&buf[0])); err != nil {
		return nil, err
	}
	return buf, nil
}

func bitsetBytes(max uint16) int {
	return int(max/8) + 1
}

func bitsetHas(bits []byte, code uint16) bool {
	index := int(code / 8)
	if index < 0 || index >= len(bits) {
		return false
	}
	mask := byte(1 << (code % 8))
	return bits[index]&mask != 0
}

func getStringIoctl(d *Device, reqFn func(uint) uint) (string, error) {
	if d == nil || d.file == nil {
		return "", ErrClosed
	}
	fd, _ := d.FD()
	buf := make([]byte, 256)
	if err := ioctl.CallPtr(fd, reqFn(uint(len(buf))), unsafe.Pointer(&buf[0])); err != nil {
		return "", err
	}
	return string(bytes.TrimRight(buf, "\x00")), nil
}

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

func readEvent(d *Device, timeout time.Duration) (Event, error) {
	if d == nil || d.file == nil {
		return Event{}, ErrClosed
	}
	fd := int(d.file.Fd())
	if err := waitReadable(fd, timeout); err != nil {
		return Event{}, err
	}
	var raw inputEvent
	buf := make([]byte, int(unsafe.Sizeof(raw)))
	if _, err := io.ReadFull(d.file, buf); err != nil {
		return Event{}, err
	}
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &raw); err != nil {
		return Event{}, err
	}
	when := time.Unix(int64(raw.Time.Sec), int64(raw.Time.Usec)*1000)
	return Event{
		When:  when,
		Kind:  EventKind(raw.Type),
		Code:  raw.Code,
		Value: raw.Value,
	}, nil
}

func writeEvent(d *Device, ev Event) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	if d.readOnly {
		return ErrReadOnly
	}
	when := ev.When
	if when.IsZero() {
		when = time.Now()
	}
	raw := inputEvent{
		Time:  syscall.NsecToTimeval(when.UnixNano()),
		Type:  uint16(ev.Kind),
		Code:  ev.Code,
		Value: ev.Value,
	}
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, raw); err != nil {
		return err
	}
	_, err := d.file.Write(buf.Bytes())
	return err
}
