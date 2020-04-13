package utils

import "testing"

func TestUnicode2UTF8(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unicode2UTF8(tt.args.text); got != tt.want {
				t.Errorf("Unicode2UTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}
