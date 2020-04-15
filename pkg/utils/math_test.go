package utils

import (
	"reflect"
	"testing"
)

func Test_Min(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "", args: args{v: []interface{}{0, 0.1, 0, 1000.87896}}, want: 0},
		{name: "", args: args{v: []interface{}{0, 0.1, 0, -1.0, -1, 1000.87896}}, want: -1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_Max(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "", args: args{v: []interface{}{0, 0.1, 0, 1000.87896}}, want: 1000.87896},
		{name: "", args: args{v: []interface{}{0, 0.1, 0, -1.0, -1, 1000.87896}}, want: 1000.87896},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkMax(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n++
		Max(0, 0.1, 0, -1.0, -1, 1000.87896)
	}
}

func TestAbs(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{"", args{0}, 0},
		{"", args{0.0}, 0.0},
		{"", args{1}, 1},
		{"", args{-1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}
