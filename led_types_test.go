package xpad

import "testing"

func TestLEDCommandAliases(t *testing.T) {
	if LEDRotate != LEDRotate1 {
		t.Fatalf("LEDRotate = %d, want %d", LEDRotate, LEDRotate1)
	}
	if LEDBlinkPrevious != LEDRotate2 {
		t.Fatalf("LEDBlinkPrevious = %d, want %d", LEDBlinkPrevious, LEDRotate2)
	}
	if LEDBlinkSlowPrevious != LEDRotate3 {
		t.Fatalf("LEDBlinkSlowPrevious = %d, want %d", LEDBlinkSlowPrevious, LEDRotate3)
	}
	if LEDRotateDual != LEDRotate4 {
		t.Fatalf("LEDRotateDual = %d, want %d", LEDRotateDual, LEDRotate4)
	}
	if LEDBlinkAllSlow != LEDBlinkFast {
		t.Fatalf("LEDBlinkAllSlow = %d, want %d", LEDBlinkAllSlow, LEDBlinkFast)
	}
	if LEDBlinkOnce != LEDBlinkSlow {
		t.Fatalf("LEDBlinkOnce = %d, want %d", LEDBlinkOnce, LEDBlinkSlow)
	}
}
