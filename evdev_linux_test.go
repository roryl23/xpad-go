//go:build linux

package xpad

import "testing"

func TestBitsetHelpers(t *testing.T) {
	cases := []struct {
		max  uint16
		want int
	}{
		{max: 0, want: 1},
		{max: 7, want: 1},
		{max: 8, want: 2},
		{max: 15, want: 2},
		{max: 16, want: 3},
	}

	for _, tc := range cases {
		if got := bitsetBytes(tc.max); got != tc.want {
			t.Fatalf("bitsetBytes(%d) = %d, want %d", tc.max, got, tc.want)
		}
	}

	bits := make([]byte, 2)
	bits[0] = 0b00010000
	bits[1] = 0b00000001

	if !bitsetHas(bits, 4) {
		t.Fatalf("bitsetHas should report set bit 4")
	}
	if bitsetHas(bits, 5) {
		t.Fatalf("bitsetHas should not report unset bit 5")
	}
	if !bitsetHas(bits, 8) {
		t.Fatalf("bitsetHas should report set bit 8")
	}
}
