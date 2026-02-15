package xpad

import "time"

// FFNewEffect requests a new effect slot from the kernel when uploading.
const FFNewEffect int16 = -1

// RumbleEffect describes a force-feedback rumble effect.
type RumbleEffect struct {
	ID     int16
	Strong uint16
	Weak   uint16
	Length time.Duration
	Delay  time.Duration
}

// NewRumbleEffect builds a rumble effect with a new effect slot.
func NewRumbleEffect(strong, weak uint16, length time.Duration) RumbleEffect {
	return RumbleEffect{ID: FFNewEffect, Strong: strong, Weak: weak, Length: length}
}
