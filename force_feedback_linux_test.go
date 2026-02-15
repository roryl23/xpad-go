//go:build linux

package xpad

import (
	"testing"
	"time"
)

func TestDurationToMillis(t *testing.T) {
	if got := durationToMillis(0); got != 0 {
		t.Fatalf("durationToMillis(0) = %d, want 0", got)
	}
	if got := durationToMillis(-1 * time.Millisecond); got != 0 {
		t.Fatalf("durationToMillis(-1ms) = %d, want 0", got)
	}
	if got := durationToMillis(42 * time.Millisecond); got != 42 {
		t.Fatalf("durationToMillis(42ms) = %d, want 42", got)
	}
	if got := durationToMillis(time.Duration(0xffff+1) * time.Millisecond); got != 0xffff {
		t.Fatalf("durationToMillis(over) = %d, want 0xffff", got)
	}
}
