package xpad

import "testing"

func TestDeviceInfoIsXpad(t *testing.T) {
	cases := []struct {
		name string
		info DeviceInfo
		want bool
	}{
		{
			name: "driver match",
			info: DeviceInfo{Driver: "xpad"},
			want: true,
		},
		{
			name: "driver substring",
			info: DeviceInfo{Driver: "hid-xpad"},
			want: true,
		},
		{
			name: "name xbox",
			info: DeviceInfo{Name: "Xbox 360 Controller"},
			want: true,
		},
		{
			name: "name x-box",
			info: DeviceInfo{Name: "X-BOX Controller"},
			want: true,
		},
		{
			name: "name unrelated",
			info: DeviceInfo{Name: "Generic Gamepad", Driver: "hid-generic"},
			want: false,
		},
	}

	for _, tc := range cases {
		if got := tc.info.IsXpad(); got != tc.want {
			t.Fatalf("%s: IsXpad() = %v, want %v", tc.name, got, tc.want)
		}
	}
}
