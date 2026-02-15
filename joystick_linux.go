//go:build linux

package xpad

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"

	"github.com/roryl23/xpad-go/internal/ioctl"
)

const joystickIOCBase = 0x6a // 'j'

func jsioCGVERSION() uint {
	return ioctl.IOR(joystickIOCBase, 0x01, ioctl.Size(uint32(0)))
}

func jsioCGAXES() uint {
	return ioctl.IOR(joystickIOCBase, 0x11, ioctl.Size(uint8(0)))
}

func jsioCGBUTTONS() uint {
	return ioctl.IOR(joystickIOCBase, 0x12, ioctl.Size(uint8(0)))
}

func jsioCGNAME(length uint) uint {
	return ioctl.IOC(ioctl.DirRead, joystickIOCBase, 0x13, length)
}

func jsioCSCORR() uint {
	return ioctl.IOW(joystickIOCBase, 0x21, ioctl.Size(JoystickCorrection{}))
}

func jsioCGCORR() uint {
	return ioctl.IOR(joystickIOCBase, 0x22, ioctl.Size(JoystickCorrection{}))
}

func jsioCSAXMAP() uint {
	return ioctl.IOW(joystickIOCBase, 0x31, ioctl.Size([AbsCnt]uint8{}))
}

func jsioCGAXMAP() uint {
	return ioctl.IOR(joystickIOCBase, 0x32, ioctl.Size([AbsCnt]uint8{}))
}

func jsioCSBTNMAP() uint {
	return ioctl.IOW(joystickIOCBase, 0x33, ioctl.Size([BtnMapLen]uint16{}))
}

func jsioCGBTNMAP() uint {
	return ioctl.IOR(joystickIOCBase, 0x34, ioctl.Size([BtnMapLen]uint16{}))
}

const (
	BtnMapLen = KeyMax - BtnMisc + 1
)

// Joystick represents an open /dev/input/js* device.
type Joystick struct {
	Path     string
	file     *os.File
	readOnly bool
}

// OpenJoystick opens a joystick device by path.
func OpenJoystick(path string) (*Joystick, error) {
	file, readOnly, err := openReadWriteOrReadOnly(path)
	if err != nil {
		return nil, err
	}
	return &Joystick{Path: path, file: file, readOnly: readOnly}, nil
}

// Close closes the joystick device.
func (j *Joystick) Close() error {
	if j == nil || j.file == nil {
		return nil
	}
	err := j.file.Close()
	j.file = nil
	j.readOnly = false
	return err
}

// FD returns the underlying file descriptor when open.
func (j *Joystick) FD() (uintptr, error) {
	if j == nil || j.file == nil {
		return 0, ErrClosed
	}
	return j.file.Fd(), nil
}

// ReadOnly reports whether the joystick handle is opened read-only.
func (j *Joystick) ReadOnly() bool {
	if j == nil {
		return true
	}
	return j.readOnly
}

// Version returns the joystick driver version.
func (j *Joystick) Version() (uint32, error) {
	if j == nil || j.file == nil {
		return 0, ErrClosed
	}
	fd, _ := j.FD()
	var value uint32
	if err := ioctl.CallPtr(fd, jsioCGVERSION(), unsafe.Pointer(&value)); err != nil {
		return 0, err
	}
	return value, nil
}

// Axes returns the number of axes.
func (j *Joystick) Axes() (uint8, error) {
	if j == nil || j.file == nil {
		return 0, ErrClosed
	}
	fd, _ := j.FD()
	var value uint8
	if err := ioctl.CallPtr(fd, jsioCGAXES(), unsafe.Pointer(&value)); err != nil {
		return 0, err
	}
	return value, nil
}

// Buttons returns the number of buttons.
func (j *Joystick) Buttons() (uint8, error) {
	if j == nil || j.file == nil {
		return 0, ErrClosed
	}
	fd, _ := j.FD()
	var value uint8
	if err := ioctl.CallPtr(fd, jsioCGBUTTONS(), unsafe.Pointer(&value)); err != nil {
		return 0, err
	}
	return value, nil
}

// Name returns the joystick name string.
func (j *Joystick) Name() (string, error) {
	if j == nil || j.file == nil {
		return "", ErrClosed
	}
	fd, _ := j.FD()
	buf := make([]byte, 128)
	if err := ioctl.CallPtr(fd, jsioCGNAME(uint(len(buf))), unsafe.Pointer(&buf[0])); err != nil {
		return "", err
	}
	return string(bytes.TrimRight(buf, "\x00")), nil
}

// AxisMap returns the joystick axis mapping.
func (j *Joystick) AxisMap() ([]uint8, error) {
	if j == nil || j.file == nil {
		return nil, ErrClosed
	}
	fd, _ := j.FD()
	var mapping [AbsCnt]uint8
	if err := ioctl.CallPtr(fd, jsioCGAXMAP(), unsafe.Pointer(&mapping[0])); err != nil {
		return nil, err
	}
	return mapping[:], nil
}

// SetAxisMap updates the joystick axis mapping.
func (j *Joystick) SetAxisMap(mapping []uint8) error {
	if j == nil || j.file == nil {
		return ErrClosed
	}
	if j.readOnly {
		return ErrReadOnly
	}
	if len(mapping) != AbsCnt {
		return fmt.Errorf("xpad: axis map length %d, want %d", len(mapping), AbsCnt)
	}
	fd, _ := j.FD()
	var arr [AbsCnt]uint8
	copy(arr[:], mapping)
	return ioctl.CallPtr(fd, jsioCSAXMAP(), unsafe.Pointer(&arr[0]))
}

// ButtonMap returns the joystick button mapping.
func (j *Joystick) ButtonMap() ([]uint16, error) {
	if j == nil || j.file == nil {
		return nil, ErrClosed
	}
	fd, _ := j.FD()
	var mapping [BtnMapLen]uint16
	if err := ioctl.CallPtr(fd, jsioCGBTNMAP(), unsafe.Pointer(&mapping[0])); err != nil {
		return nil, err
	}
	return mapping[:], nil
}

// SetButtonMap updates the joystick button mapping.
func (j *Joystick) SetButtonMap(mapping []uint16) error {
	if j == nil || j.file == nil {
		return ErrClosed
	}
	if j.readOnly {
		return ErrReadOnly
	}
	if len(mapping) != BtnMapLen {
		return fmt.Errorf("xpad: button map length %d, want %d", len(mapping), BtnMapLen)
	}
	fd, _ := j.FD()
	var arr [BtnMapLen]uint16
	copy(arr[:], mapping)
	return ioctl.CallPtr(fd, jsioCSBTNMAP(), unsafe.Pointer(&arr[0]))
}

// Correction returns the joystick correction values.
func (j *Joystick) Correction() (JoystickCorrection, error) {
	if j == nil || j.file == nil {
		return JoystickCorrection{}, ErrClosed
	}
	fd, _ := j.FD()
	var corr JoystickCorrection
	if err := ioctl.CallPtr(fd, jsioCGCORR(), unsafe.Pointer(&corr)); err != nil {
		return JoystickCorrection{}, err
	}
	return corr, nil
}

// SetCorrection updates joystick correction values.
func (j *Joystick) SetCorrection(corr JoystickCorrection) error {
	if j == nil || j.file == nil {
		return ErrClosed
	}
	if j.readOnly {
		return ErrReadOnly
	}
	fd, _ := j.FD()
	return ioctl.CallPtr(fd, jsioCSCORR(), unsafe.Pointer(&corr))
}

// ReadEvent blocks until the next joystick event or timeout.
func (j *Joystick) ReadEvent(timeout time.Duration) (JoystickEvent, error) {
	if j == nil || j.file == nil {
		return JoystickEvent{}, ErrClosed
	}
	fd := int(j.file.Fd())
	if err := waitReadable(fd, timeout); err != nil {
		return JoystickEvent{}, err
	}
	var raw jsEvent
	buf := make([]byte, int(unsafe.Sizeof(raw)))
	if _, err := io.ReadFull(j.file, buf); err != nil {
		return JoystickEvent{}, err
	}
	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, &raw); err != nil {
		return JoystickEvent{}, err
	}
	return JoystickEvent{
		Time:   raw.Time,
		Value:  raw.Value,
		Type:   JoystickEventType(raw.Type),
		Number: raw.Number,
	}, nil
}

type jsEvent struct {
	Time   uint32
	Value  int16
	Type   uint8
	Number uint8
}

// OpenJoystickDevice opens the joystick for a discovered device.
func OpenJoystickDevice(info DeviceInfo) (*Joystick, error) {
	if info.JoystickPath == "" {
		return nil, ErrNotFound
	}
	return OpenJoystick(info.JoystickPath)
}

// OpenJoystick opens the joystick for a discovered device.
func (d DeviceInfo) OpenJoystick() (*Joystick, error) {
	return OpenJoystickDevice(d)
}
