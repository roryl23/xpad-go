//go:build linux

package xpad

import (
	"time"
	"unsafe"

	"github.com/roryl23/xpad-go/internal/ioctl"
)

type ffTrigger struct {
	Button   uint16
	Interval uint16
}

type ffReplay struct {
	Length uint16
	Delay  uint16
}

type ffEnvelope struct {
	AttackLength uint16
	AttackLevel  uint16
	FadeLength   uint16
	FadeLevel    uint16
}

type ffPeriodicEffect struct {
	Waveform   uint16
	Period     uint16
	Magnitude  int16
	Offset     int16
	Phase      uint16
	Envelope   ffEnvelope
	CustomLen  uint32
	CustomData *int16
}

type ffRumbleEffect struct {
	StrongMagnitude uint16
	WeakMagnitude   uint16
}

type ffUnion struct {
	Periodic ffPeriodicEffect
}

type ffEffect struct {
	Type      uint16
	ID        int16
	Direction uint16
	Trigger   ffTrigger
	Replay    ffReplay
	Data      ffUnion
}

func (e *ffEffect) setRumble(strong, weak uint16) {
	rumble := (*ffRumbleEffect)(unsafe.Pointer(&e.Data))
	rumble.StrongMagnitude = strong
	rumble.WeakMagnitude = weak
}

// UploadRumble uploads a rumble effect and returns the assigned effect ID.
func (d *Device) UploadRumble(effect RumbleEffect) (int16, error) {
	if d == nil || d.file == nil {
		return 0, ErrClosed
	}
	if d.readOnly {
		return 0, ErrReadOnly
	}
	fd, _ := d.FD()
	ff := ffEffect{
		Type: FFRumble,
		ID:   effect.ID,
		Replay: ffReplay{
			Length: durationToMillis(effect.Length),
			Delay:  durationToMillis(effect.Delay),
		},
	}
	ff.setRumble(effect.Strong, effect.Weak)

	if err := ioctl.CallPtr(fd, evioCSFF(), unsafe.Pointer(&ff)); err != nil {
		return 0, err
	}
	return ff.ID, nil
}

// EraseEffect removes a previously uploaded effect.
func (d *Device) EraseEffect(id int16) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	if d.readOnly {
		return ErrReadOnly
	}
	fd, _ := d.FD()
	value := int32(id)
	return ioctl.CallPtr(fd, evioCRMFF(), unsafe.Pointer(&value))
}

// PlayEffect starts or updates an effect. Use repeat=0 to stop.
func (d *Device) PlayEffect(id int16, repeat int32) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	return writeEvent(d, Event{Kind: EVFF, Code: uint16(id), Value: repeat})
}

// StopEffect stops a previously uploaded effect.
func (d *Device) StopEffect(id int16) error {
	return d.PlayEffect(id, 0)
}

// SetGain sets the overall force-feedback gain (0-0xffff).
func (d *Device) SetGain(value uint16) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	return writeEvent(d, Event{Kind: EVFF, Code: FFGain, Value: int32(value)})
}

// SetAutocenter sets the force-feedback autocenter value (0-0xffff).
func (d *Device) SetAutocenter(value uint16) error {
	if d == nil || d.file == nil {
		return ErrClosed
	}
	return writeEvent(d, Event{Kind: EVFF, Code: FFAutocenter, Value: int32(value)})
}

// Rumble uploads and plays a rumble effect once.
func (d *Device) Rumble(strong, weak uint16, length time.Duration) (int16, error) {
	id, err := d.UploadRumble(NewRumbleEffect(strong, weak, length))
	if err != nil {
		return 0, err
	}
	if err := d.PlayEffect(id, 1); err != nil {
		return 0, err
	}
	return id, nil
}

func durationToMillis(d time.Duration) uint16 {
	if d <= 0 {
		return 0
	}
	ms := d / time.Millisecond
	if ms > 0xffff {
		return 0xffff
	}
	return uint16(ms)
}
